var suggestionsIndex = -1;
var suggestionTimeout = null;

function removeElement(selector) {
  var element = document.querySelector(selector);
  if (element) element.remove();
}

function closeSuggestions() {
  removeElement(".suggestions");
  document.querySelector("form").style.borderBottomLeftRadius = "8px";
  suggestionsIndex = -1;
}

function highlightSuggestion(index) {
  var suggestions = document.querySelectorAll(".suggestions li");
  suggestions.forEach(function (item) { item.className = ""; });
  if (suggestions[index]) {
    suggestions[index].className = "active";
    document.getElementById("query").value = suggestions[index].textContent;
    suggestionsIndex = index;
  }
}

function controlSuggestions(key) {
  var last = document.querySelectorAll(".suggestions li").length - 1;
  if (last < 0) return;
  if (key === "ArrowDown" || key === "ArrowRight") {
    highlightSuggestion(suggestionsIndex >= last ? 0 : suggestionsIndex + 1);
  } else if (key === "ArrowUp") {
    highlightSuggestion(suggestionsIndex <= 0 ? last : suggestionsIndex - 1);
  }
}

function openSuggestions(data) {
  removeElement(".suggestions");
  if (!data[1] || !data[1].length) return;
  var list = document.createElement("ul");
  list.className = "suggestions";
  data[1].slice(0, 8).forEach(function (suggestion) {
    var item = document.createElement("li");
    item.textContent = suggestion[0];
    item.onclick = function () {
      document.getElementById("query").value = this.textContent;
      closeSuggestions();
    };
    list.appendChild(item);
  });
  document.querySelector("form").style.borderBottomLeftRadius = "0";
  document.querySelector(".form-field").appendChild(list);
}

function setButtonError(button, message) {
  button.className = "result__error-button";
  button.textContent = message;
  button.style.backgroundColor = "#a01f15";
  button.style.color = "#ffd9d4";
  button.removeAttribute("aria-busy");
}

function setProcessing(button, message) {
  button.classList.add("is-processing");
  button.setAttribute("aria-busy", "true");
  button.innerHTML = '<span class="button-spinner" aria-hidden="true"></span><span>' + message + '</span>';
}

function createTaskId() {
  if (window.crypto && window.crypto.randomUUID) {
    return window.crypto.randomUUID();
  }
  return Date.now().toString(36) + Math.random().toString(36).slice(2);
}

function setConversionProgress(button, percentage) {
  var value = Math.max(0, Math.min(100, Math.round(percentage)));
  button.innerHTML = '<span class="button-spinner" aria-hidden="true"></span><span>A processar... ' + value + '%</span>';
}

function setSearchReady(button) {
  button.classList.remove("is-processing");
  button.removeAttribute("aria-busy");
  button.disabled = false;
  button.innerHTML = '<span class="button-label">Procurar</span>';
}

function convert(button) {
  var taskId = createTaskId();
  var params = new URLSearchParams({
    youtubelink: button.getAttribute("data-link"),
    format: button.getAttribute("data-format"),
    taskId: taskId
  });

  var socketProtocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  var socket = new WebSocket(socketProtocol + "//" + window.location.host + "/ws?taskId=" + encodeURIComponent(taskId));
  socket.onmessage = function (event) {
    try {
      var progress = JSON.parse(event.data);
      if (progress.task_id === taskId) setConversionProgress(button, progress.percentage);
    } catch (error) {
      return;
    }
  };

  fetch("/convert?" + params.toString())
    .then(function (response) {
      return response.json().then(function (body) {
        if (!response.ok || body.error) throw new Error(body.message || "Falha na conversao.");
        return body;
      });
    })
    .then(function (body) {
      button.className = "download-link";
      button.textContent = "Download";
      button.href = body.file;
      button.setAttribute("download", "");
      button.setAttribute("aria-label", "Baixar arquivo " + button.getAttribute("data-format").toUpperCase());
      button.removeAttribute("data-link");
      button.removeAttribute("aria-busy");
      button.onclick = null;
      socket.close();
    })
    .catch(function (conversionError) {
      socket.close();
      setButtonError(button, conversionError.message);
    });
}

function createDownload(video, format, link) {
  var button = document.createElement("a");
  button.href = "#";
  button.textContent = format.toUpperCase() + " Download";
  button.setAttribute("data-format", format);
  button.setAttribute("data-link", link);
  button.onclick = function (event) {
    event.preventDefault();
    if (!button.classList.contains("started")) {
      button.classList.add("started");
      setProcessing(button, "A processar...");
      convert(button);
    }
  };
  return button;
}

function createResult(video, directLink) {
  var result = document.createElement("div");
  var title = document.createElement("div");
  var info = document.createElement("div");
  var actions = document.createElement("div");
  title.className = "result__title";
  title.textContent = video.title;
  info.className = "result__infos";
  info.textContent = "YouTube - " + (video.channel || "") + " - Duracao: " + (video.duracao || "--:--:--");
  actions.className = "result__actions";
  var downloadLink = directLink || video.full_link;
  actions.appendChild(createDownload(video, "mp3", downloadLink));
  actions.appendChild(createDownload(video, "m4a", downloadLink));
  actions.appendChild(createDownload(video, "mp4", downloadLink));
  result.className = "result";
  result.style.animationDelay = (video.animationDelay || 0) + "ms";
  result.appendChild(title);
  result.appendChild(info);
  result.appendChild(actions);
  return result;
}

function showSearchError(message) {
  removeElement(".form__processing");
  removeElement(".results");
  var query = document.getElementById("query");
  query.value = "";
  query.placeholder = message;
}

function search(query) {
  if (/^(https?:\/\/)?(www\.)?(youtube\.com|youtu\.be)\//i.test(query)) {
    removeElement(".text");
    removeElement(".results");
    var directResults = document.createElement("div");
    directResults.className = "results";
    directResults.appendChild(createResult({
      title: "Video do YouTube selecionado",
      channel: "",
      duracao: ""
    }, query));
    document.body.insertBefore(directResults, document.querySelector(".footer"));
    setSearchReady(document.querySelector(".search-button"));
    return;
  }

  fetch("/search?q=" + encodeURIComponent(query))
    .then(function (response) {
      return response.json().then(function (body) {
        if (!response.ok || body.error) throw new Error(body.message || "Falha na busca.");
        return body;
      });
    })
    .then(function (body) {
      removeElement(".text");
      removeElement(".results");
      if (!body.results || !body.results.length) {
        showSearchError("Nenhum resultado encontrado.");
        return;
      }
      var results = document.createElement("div");
      results.className = "results";
      body.results.forEach(function (video, index) {
        video.animationDelay = index * 70;
        results.appendChild(createResult(video));
      });
      document.body.insertBefore(results, document.querySelector(".footer"));
    })
    .catch(function (searchError) {
      showSearchError(searchError.message);
    })
    .finally(function () {
      setSearchReady(document.querySelector(".search-button"));
    });
}

document.getElementById("query").addEventListener("keyup", function (event) {
  var value = this.value.trim();
  if (value.length < 1 || value.length > 64 || event.key === "Enter" || /https?:\/\//i.test(value)) {
    closeSuggestions();
    return;
  }
  if (["ArrowDown", "ArrowRight", "ArrowUp"].includes(event.key)) {
    controlSuggestions(event.key);
    return;
  }
  clearTimeout(suggestionTimeout);
  suggestionTimeout = setTimeout(function () {
    var script = document.createElement("script");
    script.src = "https://suggestqueries-clients6.youtube.com/complete/search?client=youtube&ds=yt&hl=pt&callback=openSuggestions&q=" + encodeURIComponent(value);
    document.body.appendChild(script);
  }, 200);
});

document.getElementById("query").addEventListener("blur", function () {
  setTimeout(closeSuggestions, 250);
});

document.querySelector("form").addEventListener("submit", function (event) {
  event.preventDefault();
  closeSuggestions();
  var query = document.getElementById("query").value.trim();
  if (query) {
    var searchButton = document.querySelector(".search-button");
    searchButton.disabled = true;
    setProcessing(searchButton, "A procurar...");
    search(query);
  }
});
