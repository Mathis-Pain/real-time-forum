// handleIncomingMessage Gérer un message entrant en temps réel
export function incomingMessage(data, currentChatUserId) {
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
// Ajouter un message EN BAS
export function appendMessage(container, sender, content, createdAt, isMine) {
  const msgEl = document.createElement('div')
  msgEl.classList.add('msg-bubble', isMine ? 'msg-sent' : 'msg-received')
  msgEl.innerHTML = `
    <span class="msg-time">${new Date(createdAt).toLocaleString()}</span>
    <strong>${sender}</strong>
    <p>${content}</p>`
  container.appendChild(msgEl)
  // Scroll automatique vers le bas pour voir le dernier message
  container.scrollTop = container.scrollHeight
}
