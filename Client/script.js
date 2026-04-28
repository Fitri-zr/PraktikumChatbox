async function sendMessage() {
  const input = document.getElementById("input");
  const chatbox = document.getElementById("chatbox");

  const text = input.value.trim();
  if (!text) return;

  chatbox.innerHTML += `<div class="message user">${text}</div>`;
  input.value = "";

  chatbox.innerHTML += `<div class="message bot" id="loading">Mengetik...</div>`;

  try {
    const res = await fetch("http://localhost:8080/chat", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ text })
    });

    const data = await res.json();

    document.getElementById("loading").remove();

    const botDiv = document.createElement("div");
botDiv.className = "message bot";
chatbox.appendChild(botDiv);

const formatted = formatAIText(data.reply);
smoothTypeHTML(botDiv, formatted);
    chatbox.scrollTop = chatbox.scrollHeight;

  } catch {
    document.getElementById("loading").remove();
    chatbox.innerHTML += `<div class="message bot">Server error</div>`;
  }
}

let chats = JSON.parse(localStorage.getItem("chats")) || [];
let currentChatId = null;

document.getElementById("input").addEventListener("keypress", function(e) {
  if (e.key === "Enter") sendMessage();
});

function smoothTypeHTML(element, text) {
  let i = 0;

  // ubah newline jadi <br>
  text = text.replace(/\n/g, "<br>");

  function typing() {
    if (i < text.length) {
      element.innerHTML += text.charAt(i);
      i++;
      setTimeout(typing, 12); // speed halus
    }
  }

  typing();
}

function formatAIText(text) {
  return text
    // bold **text**
    .replace(/\*\*(.*?)\*\*/g, "<b>$1</b>")

    // list angka
    .replace(/\n\d+\.\s/g, "<br><br>• ")

    // newline jadi paragraf
    .replace(/\n\n/g, "<br><br>")
    .replace(/\n/g, "<br>")

    // rapihin <br> berlebih
    .replace(/(<br>\s*){3,}/g, "<br><br>");
}

function smoothTypeHTML(element, html) {
  let i = 0;

  function typing() {
    if (i < html.length) {
      element.innerHTML = html.slice(0, i + 1);
      i++;
      setTimeout(typing, 10);
    }
  }

  typing();
}

function newChat() {
  const id = Date.now();

  const newChat = {
    id: id,
    title: "Chat Baru",
    messages: []
  };

  chats.push(newChat);
  currentChatId = id;

  saveChats();
  renderChatList();
  document.getElementById("chatbox").innerHTML = "";
}

function saveChats() {
  localStorage.setItem("chats", JSON.stringify(chats));
}

function renderChatList() {
  const chatList = document.getElementById("chatList");
  chatList.innerHTML = "";

  chats.forEach(chat => {
    const div = document.createElement("div");
    div.className = "chat-item";
    div.innerText = chat.title;

    div.onclick = () => loadChat(chat.id);

    chatList.appendChild(div);
  });
}

function loadChat(id) {
  currentChatId = id;

  const chat = chats.find(c => c.id === id);
  const chatbox = document.getElementById("chatbox");

  chatbox.innerHTML = "";

  chat.messages.forEach(msg => {
    chatbox.innerHTML += `<div class="message ${msg.role}">${msg.text}</div>`;
  });
}

async function sendMessage() {
  const input = document.getElementById("input");
  const chatbox = document.getElementById("chatbox");

  const text = input.value.trim();
  if (!text) return;

  // kalau belum ada chat → buat baru
  if (!currentChatId) {
    newChat();
  }

  const chat = chats.find(c => c.id === currentChatId);

  // simpan user message
  chat.messages.push({ role: "user", text });

  // 👉 set judul dari pesan pertama
  if (chat.messages.length === 1) {
    chat.title = text.slice(0, 30);
  }

  chatbox.innerHTML += `<div class="message user">${text}</div>`;
  input.value = "";

  chatbox.innerHTML += `<div class="message bot" id="loading">Mengetik...</div>`;

  try {
    const res = await fetch("http://localhost:8080/chat", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ text })
    });

    const data = await res.json();

    document.getElementById("loading").remove();

    const botDiv = document.createElement("div");
    botDiv.className = "message bot";
    chatbox.appendChild(botDiv);

    const formatted = formatAIText(data.reply);
    smoothTypeHTML(botDiv, formatted);

    // simpan bot message
    chat.messages.push({ role: "bot", text: data.reply });

    saveChats();
    renderChatList();

    chatbox.scrollTop = chatbox.scrollHeight;

  } catch {
    document.getElementById("loading").remove();
    chatbox.innerHTML += `<div class="message bot">Server error</div>`;
  }
}