// sendMessage Envoyer un message via WebSocket
export function sendMessage(ws, currentReceiverID) {
  // Récupère le champ textarea
  const textarea = document.querySelector('.message-sender')

  // Supprime les espaces inutiles avant/après
  const content = textarea.value.trim()

  // Vérification : message vide ?
  if (content === '') {
    console.warn('Message vide')
    return // On stoppe si aucun contenu
  }

  // Log pour debug
  console.log('Envoi message:', content)

  // Envoi au serveur Go via WebSocket
  // Le backend :
  // 1. Reçoit le message
  // 2. L’enregistre en base de données
  // 3. Le renvoie au destinataire
  // 4. Le renvoie à l’expéditeur (confirmation)
  // L’affichage côté front sera fait dans handleIncomingMessage()

  ws.send(
    JSON.stringify({
      type: 'message', // Type d’action
      receiver_id: currentReceiverID, // Destinataire
      content: content // Contenu du message
    })
  )

  // Reset du champ après envoi
  textarea.value = ''
}
