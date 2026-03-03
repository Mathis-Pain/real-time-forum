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

    const mainContent = document.getElementById("main-content");
    mainContent.innerHTML = `
        <div id="chat-box" style="border: 1px solid #ccc; height: 300px; overflow-y: scroll; padding: 10px; margin-bottom: 10px;">
            <p><i>Discussion avec ${userName}</i></p>
        </div>
        <form id="chat-form">
            <input type="text" id="chat-input" placeholder="Votre message..." required style="width: 80%;">
            <button type="submit">Envoyer</button>
        </form>
    `;

    const chatBox = document.getElementById("chat-box");
    const chatForm = document.getElementById("chat-form");
    const chatInput = document.getElementById("chat-input");

    // Charger l'historique
    fetch(`/messages/history?with=${encodeURIComponent(userName)}`)
        .then(res => res.json())
        .then(messages => {
            messages.reverse().forEach(m => appendMessageToChat(m));
            chatBox.scrollTop = chatBox.scrollHeight;
        });

    chatForm.onsubmit = (e) => {
        e.preventDefault();
        const content = chatInput.value.trim();
        if (!content) return;

        const msg = {
            type: "private_message",
            to: userName,
            content: content,
        };

        socket.send(JSON.stringify(msg));
        chatInput.value = "";
    };
}

function appendMessageToChat(msg) {
    const chatBox = document.getElementById("chat-box");
    if (!chatBox) return;

    const div = document.createElement("div");
    const date = new Date(msg.created_at);
    const timeStr = isNaN(date.getTime()) ? "" : date.toLocaleTimeString();

    const from = msg.from || "Moi";

    div.innerHTML = `<strong>[${timeStr}] ${from} :</strong> ${msg.content}`;
    div.style.padding = "4px 0";
    chatBox.appendChild(div);
}
