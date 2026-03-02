// Import des fonctions nécessaires au fonctionnement du chat

import {getWebSocket, addMessageHandler} from './websocket.js'
import {handleOnlineUsers} from './onlineUsers.js'
import {incomingMessage} from './incomingMessage.js'
import {openConversation} from './openConversation.js'
import {messageHistory} from './messageHistory.js'
let isLoadingHistory = false
let hasMoreMessages = true
export let currentChatUserId = null
let chatInitialized = false

// handleChatClick Fonction principale appelée lors d’un clic sur "Chat"
export function handleChatClick(userId = null, userName = null) {
  console.log('Clic chat', {userId, userName, chatInitialized})

  // Récupération du conteneur principal
  const main = document.querySelector('#main-content')

  // Récupération de la connexion WebSocket
  const ws = getWebSocket()

  // Vérifie que le WebSocket est bien connecté
  if (!ws || ws.readyState !== WebSocket.OPEN) {
    console.error('WebSocket non connecté')

    // Réessaie après 100ms si pas encore prêt
    setTimeout(() => handleChatClick(userId, userName), 100)
    return
  }

  // si on veut ouvrir directement une conversation
  if (userId && userName) {
    // Si l’interface n’est pas encore créée
    if (!chatInitialized) {
      // Création de l’interface
      initializeChatInterface(main, ws)

      // Petit délai pour s’assurer que le DOM est prêt
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
      // Si déjà initialisée, ouvrir directement
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

  // si clic simple sur "Chat" (pas de conversation précise)
  if (chatInitialized) {
    console.log('Interface déjà chargée')
    return
  }

  initializeChatInterface(main, ws)
}

// initializeChatInterface qui construit toute l’interface du chat
function initializeChatInterface(main, ws) {
  console.log('Initialisation interface chat...')

  // Injection HTML principale
  main.innerHTML = `
    <h2>Message</h2>
    <div class="messages">
      <div class="users-list"></div>
      <div style="align-items: center; font-size: 1.2rem;">
        Choisissez un utilisateur
      </div>
    </div>`

  // Mise à jour des états globaux
  chatInitialized = true
  currentChatUserId = null

  // Gestion centralisée des messages WebSocket
  addMessageHandler((data) => {
    console.log('Message WebSocket reçu:', data)

    // Mise à jour de la liste des utilisateurs en ligne
    if (data.type === 'online_users') {
      handleOnlineUsers(data.users, main, ws)
    }

    // Réception d’un nouveau message
    if (data.type === 'message') {
      incomingMessage(data, currentChatUserId)
    }

    // Réception de l’historique paginé
    if (data.type === 'message_history') {
      messageHistory(data, isLoadingHistory, hasMoreMessages)
    }
  })

  console.log('Interface chat initialisée')
}
