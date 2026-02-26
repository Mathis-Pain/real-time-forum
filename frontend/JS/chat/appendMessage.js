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
