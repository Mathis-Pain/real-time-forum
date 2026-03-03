import { openChatWith } from "./chat.js";

export function updateOnlineUsers(users) {
    const usersList = document.querySelector('.users-list');
    if (!usersList) {
        console.error('❌ .users-list introuvable pour mise à jour');
        return;
    }

    // Créer un Set des usernames en ligne
    const onlineUserNames = new Set(users);

    const userItems = usersList.querySelectorAll('.user-item');
    if (userItems.length === 0) {
        console.warn('⚠️ Aucun .user-item trouvé pour mise à jour');
        return;
    }

    userItems.forEach((userEl) => {
        const userName = userEl.dataset.userName;

        // Supprimer l'ancien indicateur s'il existe
        const existingDot = userEl.querySelector('.status-dot');
        if (existingDot) existingDot.remove();

        // Créer le point de statut
        const dot = document.createElement('span');
        dot.classList.add('status-dot');

        if (onlineUserNames.has(userName)) {
            userEl.classList.remove('offline');
            userEl.classList.add('online');
            dot.classList.add('dot-green');
        } else {
            userEl.classList.remove('online');
            userEl.classList.add('offline');
            dot.classList.add('dot-red');
        }

        userEl.prepend(dot);
    });
}
