import {getWebSocket, addMessageHandler} from './websocket.js'
import {handleOnlineUsers} from './onlineUsers.js'
import {handleIncomingMessage} from './handleIncomingMessage.js'
import {openConversation} from './openConversation.js'
import {handleMessageHistory} from './handleMessageHistory.js'

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
      setTimeout(
        () =>
          openConversation(
            main,
            ws,
            userId,
            userName,
            isLoadingHistory,
            hasMoreMessages
          ),
        100
      )
    } else {
      openConversation(
        main,
        ws,
        userId,
        userName,
        isLoadingHistory,
        hasMoreMessages
      )
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
      handleIncomingMessage(data, currentChatUserId)
    }
    if (data.type === 'message_history') {
      handleMessageHistory(data, isLoadingHistory, hasMoreMessages)
    }
  })
  console.log('Interface chat initialisée')
}
