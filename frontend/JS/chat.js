import {
  updateUsersList,
  loadPosts,
  buildMain,
  loadAllUsers,
  pendingNotifications
} from './layout.js'
let isLoading = false // Empêche les doubles requêtes
export let socket = null
let currentChatUser = null
let currentOffset = 0 // Suivi du nombre de messages chargés
let isProgrammaticScroll = false
let onlineUsers = []
// Variable timeout pour cacher l'indicateur "typing" après 2 secondes d'inactivité
let typingTimeout = null
// Horodatage pour limiter l'envoi de l'événement "typing" (ne pas spammer le serveur)
let lastTypingTime = 0
const LIMIT = 10 // Nombre de messages par "paquet"

export function getOnlineUsers() {
  return onlineUsers
}

export function initWebSocket() {
  socket = new WebSocket(`ws://${window.location.host}/ws`)

  socket.onopen = () => {
    console.log('WebSocket connected')
  }

  socket.onmessage = (event) => {
    const data = JSON.parse(event.data)
    console.log('WebSocket message received:', data)

    if (data.type === 'private_message') {
      console.log(data)
      handleIncomingPrivateMessage(data)
    }

    // Interception de l'événement WebSocket "typing" pour afficher l'indicateur
    if (data.type === 'typing') {
      handleIncomingTyping(data)
    }

    if (data.type === 'online_users') {
      onlineUsers = data.users
      updateUsersList(data.users)
    }
  }
  socket.onclose = () => {
    console.log('WebSocket closed')
  }
}

// Fonction appelée lors de la réception d'un événement "typing".
// Elle affiche l'indicateur textuel et le masque automatiquement après 2 secondes sans nouvelle frappe.
function handleIncomingTyping(msg) {
  if (currentChatUser && msg.from === currentChatUser) {
    const typingIndicator = document.getElementById('typing-indicator')
    if (typingIndicator) {
      typingIndicator.style.display = 'block'
      clearTimeout(typingTimeout)
      typingTimeout = setTimeout(() => {
        typingIndicator.style.display = 'none'
      }, 2000)
    }
  }
}

async function handleIncomingPrivateMessage(msg) {
  // 1. Ajouter la notification AVANT la reconstruction
  if (
    !currentChatUser ||
    (msg.from !== currentChatUser && msg.to !== currentChatUser)
  ) {
    pendingNotifications.add(msg.from)
  }

  // 2. Attendre que la sidebar soit reconstruite (await obligatoire)
  await loadAllUsers(msg.from)

  // 3. Afficher le message si on est dans le bon chat
  if (
    currentChatUser &&
    (msg.from === currentChatUser || msg.to === currentChatUser)
  ) {
    const chatBox = document.querySelector('.message-received')
    appendMessageToChat(msg, chatBox, false)
  }
}

export function openChatWith(userName) {
  currentChatUser = userName
  currentOffset = 0 // Reset pour la nouvelle conversation

  const main = document.getElementById('main-content')
  main.innerHTML = `
        <h2>Message avec ${userName}</h2>
        <div class="messages">
          <div class="conversation">
            <div class="message-received"></div>
            <!-- Élément d'interface "typing indicator" initialement caché, affiché temporairement lors de la réception du signal -->
            <div id="typing-indicator" style="display: none; align-self: flex-start; margin-left: 10px; color: gray; font-style: italic; font-size: 0.9rem;">
              ${userName} est en train d'écrire<span class="typing-dot">.</span><span class="typing-dot">.</span><span class="typing-dot">.</span>
            </div>
            <div class="message-content">
              <textarea class="message-sender" placeholder="Écris ton message..."></textarea>
            </div>
            <button class="send-message">Envoyer</button>
          </div>
        </div>
        <button id="returnHome">Retour à l'accueil</div>
        `
  document.getElementById('returnHome').addEventListener('click', async () => {
    const posts = await loadPosts()
    buildMain(posts)
  })
  const chatBox = main.querySelector('.message-received')
  const textarea = main.querySelector('.message-sender')
  const sendBtn = main.querySelector('.send-message')

  // 1. Premier chargement (les 10 derniers messages)
  loadHistory(userName, chatBox, true)

  chatBox.addEventListener('scroll', () => {
    if (isProgrammaticScroll) return // ← ignore les scrolls programmatiques
    if (chatBox.scrollTop <= 5 && !isLoading) {
      loadHistory(userName, chatBox, false)
    }
  })

  // Envoi de message
  sendBtn.addEventListener('click', () => sendMessage(userName, textarea))
  textarea.addEventListener('keydown', (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      sendMessage(userName, textarea)
    }
  })
  // Écouteur sur la zone de texte pour détecter la frappe de l'utilisateur.
  // Envoie un signal "typing" via WebSocket, limité à un envoi par seconde.
  textarea.addEventListener('input', () => {
    const now = Date.now()
    if (now - lastTypingTime > 1000) {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({ type: 'typing', to: userName }))
      }
      lastTypingTime = now
    }
  })
}

function loadHistory(userName, chatBox, isInitial) {
  if (isLoading) return
  isLoading = true

  fetch(
    `/messages/history?with=${encodeURIComponent(userName)}&offset=${currentOffset}&limit=${LIMIT}`
  )
    .then((res) => res.json())
    .then((messages) => {
      if (!messages || messages.length === 0) {
        isLoading = false
        return
      }
      const oldScrollHeight = chatBox.scrollHeight

      // Messages en DESC → on inverse pour ordre chronologique
      const ordered = [...messages].reverse()

      // Fragment pour insérer en une fois sans inverser l'ordre
      const fragment = document.createDocumentFragment()
      ordered.forEach((m) => fragment.appendChild(createMessageDiv(m)))
      chatBox.insertBefore(fragment, chatBox.firstChild) // ← prepend propre

      currentOffset += messages.length

      if (isInitial) {
        isProgrammaticScroll = true
        chatBox.scrollTop = chatBox.scrollHeight
        isProgrammaticScroll = false
      } else {
        isProgrammaticScroll = true
        chatBox.scrollTop = chatBox.scrollHeight - oldScrollHeight
        isProgrammaticScroll = false
      }

      isLoading = false
    })
    .catch((err) => {
      console.error('Erreur historique:', err)
      isLoading = false
    })
}

// Extraction de la création du div (utilisé aussi dans appendMessageToChat)
function createMessageDiv(msg) {
  const div = document.createElement('div')
  div.className = 'message-bubble'
  div.style.padding = '4px 0'
  const date = new Date(msg.created_at)
  const timeStr = date.toLocaleString()
  const from = msg.from || 'Moi'
  div.innerHTML = `${timeStr} <strong>${from}:</strong> `
  div.appendChild(document.createTextNode(msg.content))
  return div
}

function sendMessage(userName, textarea) {
  const content = textarea.value.trim()
  if (!content || !socket) return

  socket.send(JSON.stringify({type: 'private_message', to: userName, content}))
  textarea.value = ''

  // Re-trier la sidebar après envoi
  loadAllUsers()
}

function appendMessageToChat(msg, chatBox, isHistory) {
  if (!chatBox) return
  const div = createMessageDiv(msg)
  if (isHistory) {
    chatBox.prepend(div)
  } else {
    chatBox.appendChild(div)
    chatBox.scrollTop = chatBox.scrollHeight
  }
}
