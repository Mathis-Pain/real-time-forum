// Connexion WebSocket partagée globalement
let ws = null
let reconnectInterval = null
let messageHandlers = []

let currentUser = {id: null, name: null}

export function initWebSocket() {
  if (ws && ws.readyState === WebSocket.OPEN) {
    console.log('WebSocket déjà connecté')
    return ws
  }

  ws = new WebSocket('ws://localhost:8080/ws')

  ws.onopen = function () {
    console.log('Connecté au WebSocket')
    clearInterval(reconnectInterval)
  }

  ws.onclose = function () {
    console.log('Déconnecté du WebSocket')
    reconnectInterval = setInterval(() => {
      console.log('Tentative de reconnexion...')
      initWebSocket()
    }, 3000)
  }

  ws.onerror = function (error) {
    console.error('Erreur WebSocket:', error)
  }

  // Gestionnaire unique qui dispatch à tous les handlers
  ws.onmessage = function (event) {
    const data = JSON.parse(event.data)
    console.log('Message WebSocket reçu:', data.type)

    if (data.type === 'current_user') {
      currentUser.id = data.id
      currentUser.name = data.name
      console.log('Utilisateur identifié:', currentUser.name)
    }
    // Appeler tous les gestionnaires enregistrés
    messageHandlers.forEach((handler) => {
      try {
        handler(data)
      } catch (error) {
        console.error('Erreur dans un handler:', error)
      }
    })
  }

  return ws
}
export function getCurrentUser() {
  return currentUser
}

export function getWebSocket() {
  return ws
}

//  Ajouter un gestionnaire de messages
export function addMessageHandler(handler) {
  if (!messageHandlers.includes(handler)) {
    messageHandlers.push(handler)
    console.log('Handler ajouté, total:', messageHandlers.length)
  }
}

//  Retirer un gestionnaire
export function removeMessageHandler(handler) {
  const index = messageHandlers.indexOf(handler)
  if (index > -1) {
    messageHandlers.splice(index, 1)
    console.log('Handler retiré, total:', messageHandlers.length)
  }
}
