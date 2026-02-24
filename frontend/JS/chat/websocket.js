// Connexion WebSocket partagée globalement (Singleton)

// Instance unique du WebSocket
let ws = null

// Intervalle utilisé pour la reconnexion automatique
let reconnectInterval = null

// Liste des fonctions qui écoutent les messages entrants
let messageHandlers = []

// Objet représentant l’utilisateur actuellement connecté
let currentUser = {id: null, name: null}

// initWebSocket Initialiser la connexion WebSocket
export function initWebSocket() {
  // Si déjà connecté, on retourne l’instance existante
  if (ws && ws.readyState === WebSocket.OPEN) {
    console.log('WebSocket déjà connecté')
    return ws
  }

  // Création d’une nouvelle connexion
  ws = new WebSocket('ws://localhost:8080/ws')

  // Événement : connexion ouverte
  ws.onopen = function () {
    console.log('Connecté au WebSocket')

    // Stoppe la tentative de reconnexion si active
    clearInterval(reconnectInterval)
  }

  // Événement : connexion fermée
  ws.onclose = function () {
    console.log('Déconnecté du WebSocket')

    // Tentative de reconnexion toutes les 3 secondes
    reconnectInterval = setInterval(() => {
      console.log('Tentative de reconnexion...')
      initWebSocket()
    }, 3000)
  }

  // Événement : erreur WebSocket
  ws.onerror = function (error) {
    console.error('Erreur WebSocket:', error)
  }

  // Événement : message reçu
  ws.onmessage = function (event) {
    // Conversion du message JSON reçu
    const data = JSON.parse(event.data)

    console.log('Message WebSocket reçu:', data.type)

    // Si le serveur envoie l'utilisateur courant
    if (data.type === 'current_user') {
      currentUser.id = data.id
      currentUser.name = data.name

      console.log('Utilisateur identifié:', currentUser.name)
    }

    // Dispatcher : appel de tous les handlers enregistrés
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

// Retourner l'utilisateur courant
export function getCurrentUser() {
  return currentUser
}

// Retourner l'instance WebSocket
export function getWebSocket() {
  return ws
}

// Ajouter un gestionnaire de messages
export function addMessageHandler(handler) {
  // Évite les doublons
  if (!messageHandlers.includes(handler)) {
    messageHandlers.push(handler)
    console.log('Handler ajouté, total:', messageHandlers.length)
  }
}

// Retirer un gestionnaire
export function removeMessageHandler(handler) {
  const index = messageHandlers.indexOf(handler)

  if (index > -1) {
    messageHandlers.splice(index, 1)
    console.log('Handler retiré, total:', messageHandlers.length)
  }
}
