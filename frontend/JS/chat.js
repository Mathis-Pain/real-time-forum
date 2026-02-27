import {handleMessageHistory} from './chat/handleMessageHistory.js'
// frontend/JS/chat.js

export function initChat(receiverNickname) {
  const mainContent = document.getElementById('main-content')

  const chatBox = document.getElementById('chat-box')
  const chatForm = document.getElementById('chat-form')
  const chatInput = document.getElementById('chat-input')

  // 2. Connexion au WebSocket (on passe le receiver dans l'URL comme prévu dans ton Go)
  const socket = new WebSocket(
    `ws://${window.location.host}/ws?receiver=${receiverNickname}`
  )

  // 3. Réception des messages (ce que ton Go envoie via WriteMessagesFromBddToUserScreen)
  socket.onmessage = function (event) {
    const messageElement = document.createElement('div')
    messageElement.textContent = event.data // Le contenu brut envoyé par Go
    messageElement.style.padding = '5px'
    messageElement.style.borderBottom = '1px dotted #eee'
    chatBox.appendChild(messageElement)
    chatBox.scrollTop = chatBox.scrollHeight // Scroll auto vers le bas
  }

  // 4. Envoi de message
  chatForm.onsubmit = function (e) {
    e.preventDefault()
    if (chatInput.value.trim() !== '') {
      socket.send(chatInput.value) // Envoi au format texte (MT=1)

      // On l'affiche aussi localement pour nous
      const myMsg = document.createElement('div')
      myMsg.innerHTML = `<strong>Moi:</strong> ${chatInput.value}`
      chatBox.appendChild(myMsg)
      chatInput.value = ''
      handleMessageHistory()
    }
  }

  socket.onclose = function () {
    console.log('Connexion WebSocket fermée.')
  }
}
