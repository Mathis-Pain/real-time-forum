function Chat() {
  const ws = new WebSocket('ws://localhost:8080/ws')

  ws.onopen = () => {
    console.log('connexion etablie')
  }
  ws.send = () => {
    console.log('envoi de message')
  }
  ws.onmessage = function (e) {
    console.log('message brut reçu :', e.data)
    console.log(e)
    const data = JSON.parse(e.data)
    console.log(data)
    if (data.type == 'message') {
      console.log(data)
    }
    if (data.type == 'online-users') {
      console.log(data)
    }
  }
  ws.onerror = (error) => {
    console.log('erreur de connexion', error)
  }
  ws.onclose = () => {
    console.log('connexion perdu')
  }
}
