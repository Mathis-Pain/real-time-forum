let socket = null;
let currentChatUser = null;

export function initWebSocket() {
    socket = new WebSocket(`ws://${window.location.host}/ws`);

    socket.onopen = () => {
        console.log("WebSocket connected");
    };

    socket.onmessage = (event) => {
        const data = JSON.parse(event.data);

        if (data.type === "private_message") {
            handleIncomingPrivateMessage(data);
        }
    };

    socket.onclose = () => {
        console.log("WebSocket closed");
    };
}

function handleIncomingPrivateMessage(msg) {
    if (currentChatUser && (msg.from === currentChatUser || msg.to === currentChatUser)) {
        appendMessageToChat(msg);
    }
}

export function openChatWith(userName) {
    currentChatUser = userName;

    const main = document.getElementById("main-content");
    main.innerHTML = `
        <h2>Message avec ${userName}</h2>
        <div class="messages">
          <div class="users-list"></div>
          <div class="conversation">
            <div class="message-received"></div>
            <div class="message-content">
              <textarea class="message-sender" placeholder="Écris ton message..."></textarea>
            </div>
            <button class="send-message">Envoyer</button>
          </div>
        </div>`;

    const chatBox = main.querySelector(".message-received");
    const textarea = main.querySelector(".message-sender");
    const sendBtn = main.querySelector(".send-message");

    // Charger l'historique
    fetch(`/messages/history?with=${encodeURIComponent(userName)}`)
        .then(res => res.json())
        .then(messages => {
            if (!messages) return;
            messages.reverse().forEach(m => appendMessageToChat(m));
            chatBox.scrollTop = chatBox.scrollHeight;
        })
        .catch(err => console.error("Erreur chargement historique:", err));

    // Envoi avec le bouton
    sendBtn.addEventListener("click", () => sendMessage(userName, textarea));

    // Envoi avec Entrée (Shift+Entrée pour saut de ligne)
    textarea.addEventListener("keydown", (e) => {
        if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault();
            sendMessage(userName, textarea);
        }
    });
}

function sendMessage(userName, textarea) {
    const content = textarea.value.trim();
    if (!content || !socket) return;

    const msg = {
        type: "private_message",
        to: userName,
        content: content,
    };

    socket.send(JSON.stringify(msg));
    textarea.value = "";
}

function appendMessageToChat(msg) {
    const chatBox = document.querySelector(".message-received");
    if (!chatBox) return;

    const div = document.createElement("div");
    const date = new Date(msg.created_at);
    const timeStr = date.toLocaleDateString() + " " + date.toLocaleTimeString()
    const from = msg.from || "Moi";

    div.classList.add("message-bubble");
    div.innerHTML = `<strong>[${timeStr}] ${from} :</strong> ${msg.content}`;
    div.style.padding = "4px 0";
    chatBox.appendChild(div);
}

