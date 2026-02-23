// Envoyer un message
export function sendMessage(ws, currentReceiverID) {
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
