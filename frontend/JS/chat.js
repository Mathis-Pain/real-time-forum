// frontend/JS/chat.js

export function initChat(receiverNickname) {
    const mainContent = document.getElementById("main-content");

    // 1. On crée une interface très simple
    mainContent.innerHTML = `
        <div id="chat-box" style="border: 1px solid #ccc; height: 300px; overflow-y: scroll; padding: 10px; margin-bottom: 10px;">
            <p><i>Début de la discussion avec ${receiverNickname}</i></p>
        </div>
        <form id="chat-form">
            <input type="text" id="chat-input" placeholder="Votre message..." required style="width: 80%;">
            <button type="submit">Envoyer</button>
        </form>
    `;

    const chatBox = document.getElementById("chat-box");
    const chatForm = document.getElementById("chat-form");
    const chatInput = document.getElementById("chat-input");

    // 2. Connexion au WebSocket (on passe le receiver dans l'URL comme prévu dans ton Go)
    const socket = new WebSocket(`ws://${window.location.host}/ws?receiver=${receiverNickname}`);

    // 3. Réception des messages (ce que ton Go envoie via WriteMessagesFromBddToUserScreen)
    socket.onmessage = function(event) {
        const messageElement = document.createElement("div");
        messageElement.textContent = event.data; // Le contenu brut envoyé par Go
        messageElement.style.padding = "5px";
        messageElement.style.borderBottom = "1px dotted #eee";
        chatBox.appendChild(messageElement);
        chatBox.scrollTop = chatBox.scrollHeight; // Scroll auto vers le bas
    };

    // 4. Envoi de message
    chatForm.onsubmit = function(e) {
        e.preventDefault();
        if (chatInput.value.trim() !== "") {
            socket.send(chatInput.value); // Envoi au format texte (MT=1)
            
            // On l'affiche aussi localement pour nous
            const myMsg = document.createElement("div");
            myMsg.innerHTML = `<strong>Moi:</strong> ${chatInput.value}`;
            myMsg.style.color = "blue";
            chatBox.appendChild(myMsg);
            
            chatInput.value = "";
        }
    };

    socket.onclose = function() {
        console.log("Connexion WebSocket fermée.");
    };
}