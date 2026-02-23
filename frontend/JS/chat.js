import {getWebSocket, addMessageHandler} from './websocket.js'
import {handleOnlineUsers} from './onlineUsers.js'

let currentOffset = 0
let currentReceiverID = null
let isLoadingHistory = false
let hasMoreMessages = true
let currentChatUserId = null
let chatInitialized = false

// Fonction principale pour gérer le chat
export function handleChatClick(e, userId = null, userName = null) {
  if (e) e.preventDefault()
  console.log('Clic chat', {userId, userName, chatInitialized})

  const main = document.querySelector('#main-content')
  const ws = getWebSocket()

  if (!ws || ws.readyState !== WebSocket.OPEN) {
    console.error('WebSocket non connecté')
    setTimeout(() => handleChatClick(null, userId, userName), 100)
    return
  }

  // Si userId et userName sont fournis, ouvrir directement la conversation
  if (userId && userName) {
    if (!chatInitialized) {
      initializeChatInterface(main, ws)
      setTimeout(() => openConversation(main, ws, userId, userName), 100)
    } else {
      openConversation(main, ws, userId, userName)
    }
    return
  }

  if (chatInitialized) {
    console.log('Interface déjà chargée')
    return
  }

  initializeChatInterface(main, ws)
}

// Initialiser l'interface complète du chat
function initializeChatInterface(main, ws) {
  console.log('Initialisation interface chat...')

  main.innerHTML = `
    <h2>Message</h2>
    <div class="messages">
      <div class="users-list"></div>
      <div style="align-items: center; font-size: 1.2rem;">Choisissez un utilisateur</div>
    </div>`

  chatInitialized = true
  currentChatUserId = null

  addMessageHandler((data) => {
    console.log('Message WebSocket reçu:', data)

    if (data.type === 'online_users') {
      handleOnlineUsers(data.users, main, ws)
    }

    if (data.type === 'message') {
      handleIncomingMessage(data)
    }

    if (data.type === 'message_history') {
      handleMessageHistory(data)
    }
  })

  console.log('Interface chat initialisée')
}

// Gérer un message entrant en temps réel
function handleIncomingMessage(data) {
  const senderId = data.sender_id

  if (currentChatUserId === senderId || data.is_mine) {
    // Message de la conversation active (envoyé ou reçu)
    const receivedDiv = document.querySelector('.message-received')
    if (receivedDiv) {
      appendMessage(
        receivedDiv,
        data.sender,
        data.content,
        data.created_at,
        data.is_mine
      )
    }
  } else {
    // Notification sur l'utilisateur dans la sidebar
    const userItem = document.querySelector(
      `.user-item[data-user-id="${senderId}"]`
    )
    if (userItem) {
      userItem.classList.add('has-notification')
    }
  }
}

// Gérer l'historique des messages (avec pagination)
function handleMessageHistory(data) {
  const receivedDiv = document.querySelector('.message-received')
  if (!receivedDiv) {
    console.warn('.message-received introuvable')
    return
  }

  isLoadingHistory = false
  hasMoreMessages = data.has_more

  if (data.offset === 0) {
    receivedDiv.innerHTML = ''
    if (data.messages) {
      data.messages.forEach((msg) => {
        appendMessage(
          receivedDiv,
          msg.sender,
          msg.content,
          msg.created_at,
          msg.is_mine
        )
      })
    }
    receivedDiv.scrollTop = receivedDiv.scrollHeight
  } else {
    const previousHeight = receivedDiv.scrollHeight

    if (data.messages) {
      data.messages.reverse().forEach((msg) => {
        prependMessage(
          receivedDiv,
          msg.sender,
          msg.content,
          msg.created_at,
          msg.is_mine
        )
      })
    }

    receivedDiv.scrollTop = receivedDiv.scrollHeight - previousHeight
  }

  console.log(
    `Historique chargé (offset: ${data.offset}, has_more: ${data.has_more})`
  )
}

// Ouvrir une conversation
function openConversation(main, ws, userId, userName) {
  console.log(`Ouverture conversation: ${userName} (ID: ${userId})`)

  currentOffset = 0
  currentReceiverID = userId
  currentChatUserId = userId
  hasMoreMessages = true
  isLoadingHistory = false

  main.innerHTML = `
    <h2>Message avec ${userName}</h2>
    <div class="messages">
      <div class="users-list"></div>
      <div class="conversation">
        <div class="message-received"></div>
        <div class="message-content">
          <textarea class="message-sender" placeholder="Écris ton message..."></textarea>
          </div>
          <button class="send-message">Envoyer</button>
      </div>
    </div>`

  // Charger les 10 premiers messages
  ws.send(
    JSON.stringify({
      type: 'get_history',
      receiver_id: userId,
      offset: 0
    })
  )

  currentOffset = 10

  // Scroll pour charger plus de messages
  const receivedDiv = document.querySelector('.message-received')
  receivedDiv.addEventListener('scroll', () => {
    if (receivedDiv.scrollTop === 0 && !isLoadingHistory && hasMoreMessages) {
      console.log('Chargement messages plus anciens...')
      isLoadingHistory = true
      ws.send(
        JSON.stringify({
          type: 'get_history',
          receiver_id: currentReceiverID,
          offset: currentOffset
        })
      )
      currentOffset += 10
    }
  })

  // Envoyer (bouton)
  document.querySelector('.send-message').addEventListener('click', () => {
    sendMessage(ws)
  })

  // Envoyer (touche Entrée)
  document
    .querySelector('.message-sender')
    .addEventListener('keypress', (e) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault()
        sendMessage(ws)
      }
    })
}

// Envoyer un message
function sendMessage(ws) {
  const textarea = document.querySelector('.message-sender')
  const content = textarea.value.trim()

  if (content === '') {
    console.warn('Message vide')
    return
  }

  console.log('Envoi message:', content)

  // ✅ Go reçoit, sauvegarde et renvoie le message
  // handleIncomingMessage() se charge de l'affichage
  ws.send(
    JSON.stringify({
      type: 'message',
      receiver_id: currentReceiverID,
      content: content
    })
  )

  textarea.value = ''
}

// Ajouter un message EN BAS
function appendMessage(container, sender, content, createdAt, isMine) {
  const msgEl = document.createElement('div')
  msgEl.classList.add('msg-bubble', isMine ? 'msg-sent' : 'msg-received')
  msgEl.innerHTML = `
    <span class="msg-time">${new Date(createdAt).toLocaleTimeString()}</span>
    <strong>${sender}</strong>
    <p>${content}</p>`
  container.appendChild(msgEl)
  container.scrollTop = container.scrollHeight
}

// Ajouter un message EN HAUT (pagination)
function prependMessage(container, sender, content, createdAt, isMine) {
  const msgEl = document.createElement('div')
  msgEl.classList.add('msg-bubble', isMine ? 'msg-sent' : 'msg-received')
  msgEl.innerHTML = `
    <span class="msg-time">${new Date(createdAt).toLocaleTimeString()}</span>
    <strong>${sender}</strong>
    <p>${content}</p>`
  container.prepend(msgEl)
}
