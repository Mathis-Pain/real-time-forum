// Import de la fonction qui ajoute un message en bas
import {appendMessage} from './appendMessage.js'

// handleMessageHistory Gérer l’historique des messages (avec pagination)
export function handleMessageHistory(data, isLoadingHistory, hasMoreMessages) {
  // Récupération du conteneur des messages
  const receivedDiv = document.querySelector('.message-received')

  // Sécurité : si le conteneur n’existe pas on stoppe
  if (!receivedDiv) {
    console.warn('.message-received introuvable')
    return
  }

  // On indique que le chargement est terminé
  isLoadingHistory = false

  // Mise à jour de l’état indiquant s’il reste des messages à charger
  hasMoreMessages = data.has_more

  //  Chargement initial (offset = 0) offset représente combien d’éléments on a déjà chargés.
  if (data.offset === 0) {
    // On vide la conversation
    receivedDiv.innerHTML = ''

    // Si des messages sont présents
    if (data.messages) {
      data.messages.forEach((msg) => {
        // On les ajoute EN BAS (ordre chronologique)
        appendMessage(
          receivedDiv,
          msg.sender,
          msg.content,
          msg.created_at,
          msg.is_mine
        )
      })
    }

    // Scroll automatique vers le bas
    receivedDiv.scrollTop = receivedDiv.scrollHeight
  }

  // Si Pagination (offset > 0)
  else {
    // On sauvegarde la hauteur AVANT ajout
    const previousHeight = receivedDiv.scrollHeight

    if (data.messages) {
      data.messages.reverse().forEach((msg) => {
        // On ajoute les anciens messages EN HAUT
        prependMessage(
          receivedDiv,
          msg.sender,
          msg.content,
          msg.created_at,
          msg.is_mine
        )
      })
    }

    // Ajustement du scroll pour éviter le "saut visuel"
    receivedDiv.scrollTop = receivedDiv.scrollHeight - previousHeight
  }

  console.log(
    `Historique chargé (offset: ${data.offset}, has_more: ${data.has_more})`
  )
}

// Ajouter un message EN HAUT (utilisé pour la pagination)
function prependMessage(container, sender, content, createdAt, isMine) {
  // Création de la bulle
  const msgEl = document.createElement('div')

  // Style différent selon l’expéditeur
  msgEl.classList.add('msg-bubble', isMine ? 'msg-sent' : 'msg-received')

  // Contenu HTML du message
  msgEl.innerHTML = `
    <span class="msg-time">
      ${new Date(createdAt).toLocaleString()}
    </span>
    <strong>${sender}</strong>
    <p>${content}</p>
  `

  // Insertion en haut du conteneur
  container.prepend(msgEl)
}
