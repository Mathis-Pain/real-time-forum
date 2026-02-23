import {appendMessage} from './appendMessage.js'

// Gérer l'historique des messages (avec pagination)
export function handleMessageHistory(data, isLoadingHistory, hasMoreMessages) {
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
      data.messages.forEach((msg) => {
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
