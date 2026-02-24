import {appendMessage} from './appendMessage.js'

// handleIncomingMessage Gérer un message entrant en temps réel
export function handleIncomingMessage(data, currentChatUserId) {
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
