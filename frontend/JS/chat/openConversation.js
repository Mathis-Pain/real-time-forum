import {sendMessage} from './sendMessage.js'

// Ouvrir une conversation
export function openConversation(
  main,
  ws,
  userId,
  userName,
  isLoadingHistory,
  hasMoreMessages
) {
  console.log(`Ouverture conversation: ${userName} (ID: ${userId})`)

  let currentOffset = 0
  let currentReceiverID = userId

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
    sendMessage(ws, currentReceiverID)
  })

  // Envoyer (touche Entrée)
  document
    .querySelector('.message-sender')
    .addEventListener('keypress', (e) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault()
        sendMessage(ws, currentReceiverID)
      }
    })
}
