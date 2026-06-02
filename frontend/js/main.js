// Инициализация интерфейса после загрузки DOM
document.addEventListener("DOMContentLoaded", () => {
    
    // 1. СЕЛЕКТОРЫ И ГЛОБАЛЬНОЕ СОСТОЯНИЕ
    const els = {
        containersList: document.getElementById("containers-list"),
        controlPanel: document.getElementById("control-panel"),
        panelTitle: document.getElementById("selected-container-title"),
        btnCreate: document.getElementById("btn-create"),
        btnLogout: document.getElementById("btn-logout"),
        
        // Вкладки навигации (Tabs)
        tabContainer: document.getElementById("btn-tab-container"),
        tabLogs: document.getElementById("btn-tab-logs"),
        tabConsole: document.getElementById("btn-tab-console"),
        
        // Окна вкладок (Panes)
        paneContainer: document.getElementById("pane-container"),
        paneLogs: document.getElementById("pane-logs"),
        paneConsole: document.getElementById("pane-console"),

        // Управление контейнером (Вкладка Container)
        statusBadge: document.getElementById("container-status"),
        btnToggle: document.getElementById("btn-action-toggle"),
        btnRestart: document.getElementById("btn-action-restart"),
        btnDelete: document.getElementById("btn-action-delete"),

        // Компоненты логов и консоли
        logsOutput: document.getElementById("logs-output"),
        btnRefreshLogs: document.getElementById("btn-refresh-logs"),
        inputConsole: document.getElementById("input-console"),
        consoleOutput: document.getElementById("console-output"),

        // Модальное окно: Создание
        modalCreate: document.getElementById("modal-create"),
        btnModalClose: document.getElementById("btn-modal-close"),
        btnModalSubmit: document.getElementById("btn-modal-submit"),
        inputName: document.getElementById("input-name"),
        inputUrl: document.getElementById("input-url"),
        modalError: document.getElementById("modal-error"),

        // Модальное окно: Удаление
        modalDelete: document.getElementById("modal-delete"),
        btnDeleteClose: document.getElementById("btn-delete-close"),
        btnDeleteCancel: document.getElementById("btn-delete-cancel"),
        btnDeleteConfirm: document.getElementById("btn-delete-confirm"),
        deleteContainerName: document.getElementById("delete-container-name"),
        modalDeleteError: document.getElementById("modal-delete-error")
    };

    let selectedContainerId = null;
    let selectedContainerName = "";
    let currentStatus = ""; // "running" или "stopped"

    // 2. ВСПОМОГАТЕЛЬНЫЕ УТИЛИТЫ (ПОДФУНКЦИИ)
    const toggleVisibility = (el, isVisible) => {
        if (!el) return;
        el.classList.toggle("hidden", !isVisible);
        el.setAttribute("aria-hidden", (!isVisible).toString());
    };

    const scrollTerminalToBottom = (pane) => {
        if (!pane) return;
        const viewer = pane.querySelector(".terminal-viewer");
        if (viewer) viewer.scrollTop = viewer.scrollHeight;
    };

        // 3. УПРАВЛЕНИЕ ВКЛАДКАМИ (TABS)
    const tabsConfig = [
        { tab: els.tabContainer, pane: els.paneContainer, onOpen: handleContainerTabOpen },
        { tab: els.tabLogs, pane: els.paneLogs, onOpen: handleLogsTabOpen },
        { tab: els.tabConsole, pane: els.paneConsole, onOpen: handleConsoleTabOpen }
    ];

    function initTabs() {
        tabsConfig.forEach(({ tab }) => {
            if (tab) tab.addEventListener("click", () => switchTab(tab));
        });
    }

    function switchTab(targetTab) {
        tabsConfig.forEach(({ tab, pane, onOpen }) => {
            const isActive = tab === targetTab;
            if (tab) tab.classList.toggle("active", isActive);
            toggleVisibility(pane, isActive);
            if (isActive && onOpen) onOpen();
        });
    }

    // Подфункции жизненного цикла вкладок
    function handleContainerTabOpen() {
        updateContainerStatusUI();
    }

    async function handleLogsTabOpen() {
        if (!els.logsOutput) return;
        els.logsOutput.textContent = "Получение логов...";
        try {
            const res = await loadData("POST", `/containers/logs/`, { Id: selectedContainerId, Tail: 100 });
            els.logsOutput.textContent = res.logs || "Логи пусты";
            scrollTerminalToBottom(els.paneLogs);
        } catch (err) {
            els.logsOutput.textContent = `Ошибка загрузки логов: ${err.message}`;
        }
    }

    function handleConsoleTabOpen() {
        if (!els.consoleOutput || !els.inputConsole) return;
        els.consoleOutput.innerHTML = "<div>Подключено к терминалу. Готов к вводу команд.</div>";
        els.inputConsole.focus();
    }

    // 4. СПИСОК КОНТЕЙНЕРОВ
    async function initContainersList() {
        try {
            const res = await loadData("GET", "/containers/get");
            if (!els.containersList) return;
            els.containersList.innerHTML = "";

            if (!res?.containers_list?.length) {
                els.containersList.innerHTML = '<li class="loading">Контейнеры не найдены</li>';
                return;
            }

            renderContainers(res.containers_list);
        } catch (err) {
            if (els.containersList) {
                els.containersList.innerHTML = `<li class="loading" style="color:var(--danger)">Ошибка: ${err.message}</li>`;
            }
        }
    }

    function renderContainers(containers) {
        const fragment = document.createDocumentFragment();

        containers.forEach(({ Id, Name, Status }) => {
            const li = document.createElement("li");
            li.className = "container-item";
            li.textContent = Name;
            li.dataset.id = Id;

            if (Id === selectedContainerId) li.classList.add("active");

            li.addEventListener("click", () => selectContainer(li, Id, Name, Status));
            fragment.appendChild(li);
        });

        els.containersList.appendChild(fragment);
    }

    function selectContainer(element, id, name, status) {
        document.querySelectorAll(".container-item").forEach(el => el.classList.remove("active"));
        element.classList.add("active");

        selectedContainerId = id;
        selectedContainerName = name;
        currentStatus = status === "running" ? "running" : "stopped";

        if (els.panelTitle) els.panelTitle.textContent = `Управление: ${name}`;
        
        toggleVisibility(els.controlPanel, true);
        switchTab(els.tabContainer);
    }

    // 5. ДЕЙСТВИЯ НАД КОНТЕЙНЕРОМ (START / STOP / RESTART)
    function initContainerActions() {
        if (els.btnToggle) els.btnToggle.addEventListener("click", toggleContainerState);
        if (els.btnRestart) els.btnRestart.disabled = false; // Сбрасываем блокировку при инициализации
        if (els.btnRestart) els.btnRestart.addEventListener("click", restartContainer);
        if (els.btnRefreshLogs) els.btnRefreshLogs.addEventListener("click", handleLogsTabOpen);
        if (els.inputConsole) els.inputConsole.addEventListener("keydown", handleConsoleCommand);
    }

    function updateContainerStatusUI() {
        if (!els.statusBadge || !els.btnToggle) return;

        els.statusBadge.className = `status-badge ${currentStatus}`;
        els.statusBadge.textContent = currentStatus === "running" ? "Запущен" : "Остановлен";

        if (currentStatus === "running") {
            els.btnToggle.textContent = "Остановить";
            els.btnToggle.className = "btn btn-danger";
        } else {
            els.btnToggle.textContent = "Запустить";
            els.btnToggle.className = "btn btn-success";
        }
    }

    async function toggleContainerState() {
        const action = currentStatus === "running" ? "stop" : "start";
        try {
            els.btnToggle.disabled = true;
            await loadData("POST", `/containers/${action}/`, {Id: selectedContainerId});
            currentStatus = action === "start" ? "running" : "stopped";
            updateContainerStatusUI();
            await initContainersList();
        } catch (err) {
            alert(`Ошибка изменения состояния: ${err.message}`);
        } finally {
            els.btnToggle.dsisabled = false;
        }
    }

    async function restartContainer() {
        try {
            els.btnRestart.disabled = true;
            // ИСПРАВЛЕНО: отправка POST /containers/restart/:id без тела
            await loadData("POST", `/containers/restart/`, { Id: selectedContainerId});
            currentStatus = "running";
            updateContainerStatusUI();
            await initContainersList();
            alert("Контейнер успешно перезапущен");
        } catch (err) {
            alert(`Ошибка рестарта: ${err.message}`);
        } finally {
            els.btnRestart.disabled = false;
        }
    }

    async function handleConsoleCommand(e) {
        if (e.key !== "Enter") return;
        const command = els.inputConsole.value.trim();
        if (!command) return;

        els.consoleOutput.innerHTML += `<div><span style="color:var(--success)">#</span> ${command}</div>`;
        els.inputConsole.value = "";

        try {
            // Оставляем как есть, либо настройте под свой exec роут, если он появится
            const res = await loadData("POST", "/containers/exec", { Id: selectedContainerId, cmd: command });
            els.consoleOutput.innerHTML += `<div>${res.output || "команда выполнена"}</div>`;
        } catch (err) {
            els.consoleOutput.innerHTML += `<div style="color:var(--danger)">Ошибка: ${err.message}</div>`;
        }
        scrollTerminalToBottom(els.paneConsole);
    }
    
    // 6. МОДАЛЬНОЕ ОКНО: СОЗДАНИЕ
    function initCreateModal() {
        if (els.btnCreate) els.btnCreate.addEventListener("click", openCreateModal);
        if (els.btnModalClose) els.btnModalClose.addEventListener("click", closeCreateModal);
        if (els.btnModalSubmit) els.btnModalSubmit.addEventListener("click", submitCreateContainer);
        
        if (els.modalCreate) {
            els.modalCreate.addEventListener("click", (e) => {
                if (e.target === els.modalCreate) closeCreateModal();
            });
        }

        [els.inputName, els.inputUrl].forEach(input => {
            if (input) input.addEventListener("keydown", (e) => {
                if (e.key === "Enter") { e.preventDefault(); submitCreateContainer(); }
            });
        });
    }

    function openCreateModal() {
        if (els.inputName) els.inputName.value = "";
        if (els.inputUrl) els.inputUrl.value = "";
        toggleVisibility(els.modalError, false);
        toggleVisibility(els.modalCreate, true);
        if (els.inputName) els.inputName.focus();
    }

    function closeCreateModal() {
        toggleVisibility(els.modalCreate, false);
    }

    async function submitCreateContainer() {
        if (!els.inputName || !els.inputUrl) return;
        const name = els.inputName.value.trim();
        const url = els.inputUrl.value.trim();

        if (!name) return showCreateModalError("Поле Название не может быть пустым");
        if (!url) return showCreateModalError("Поле URL не может быть пустым");

        try {
            els.btnModalSubmit.disabled = true;
            await loadData("POST", "/containers/create/", { name, url });
            closeCreateModal();
            await initContainersList();
        } catch (err) {
            showCreateModalError(err.message);
        } finally {
            els.btnModalSubmit.disabled = false;
        }
    }

    function showCreateModalError(message) {
        if (!els.modalError) return;
        els.modalError.textContent = message;
        toggleVisibility(els.modalError, true);
    }

    // 7. МОДАЛЬНОЕ ОКНО: УДАЛЕНИЕ
    function initDeleteModal() {
        if (els.btnDelete) els.btnDelete.addEventListener("click", openDeleteModal);
        if (els.btnDeleteConfirm) els.btnDeleteConfirm.addEventListener("click", submitDeleteContainer);
        
        [els.btnDeleteClose, els.btnDeleteCancel].forEach(btn => {
            if (btn) btn.addEventListener("click", closeDeleteModal);
        });
    }

    function openDeleteModal() {
        if (els.deleteContainerName) els.deleteContainerName.textContent = selectedContainerName;
        toggleVisibility(els.modalDeleteError, false);
        toggleVisibility(els.modalDelete, true);
    }

    function closeDeleteModal() {
        toggleVisibility(els.modalDelete, false);
    }

    async function submitDeleteContainer() {
        try {
            els.btnDeleteConfirm.disabled = true;
            toggleVisibility(els.modalDeleteError, false);

            // ИСПРАВЛЕНО под роут DELETE /containers/delete/:id
            await loadData("DELETE", `/containers/delete/`, { Id: selectedContainerId});
            
            closeDeleteModal();
            toggleVisibility(els.controlPanel, false);
            selectedContainerId = null;
            selectedContainerName = "";
            await initContainersList();
        } catch (err) {
            if (els.modalDeleteError) {
                els.modalDeleteError.textContent = err.message || "Не удалось удалить контейнер";
                toggleVisibility(els.modalDeleteError, true);
            }
        } finally {
            els.btnDeleteConfirm.disabled = false;
        }
    }

    // 8. СИСТЕМНЫЕ ОБРАБОТЧИКИ И СТАРТ
    function initGlobalEvents() {
        document.addEventListener("keydown", (e) => {
            if (e.key === "Escape") {
                closeCreateModal();
                closeDeleteModal();
            }
        });

        if (els.btnLogout) {
            els.btnLogout.addEventListener("click", async () => {
                try {
                    await loadData("POST", "/logout");
                } catch (err) {
                    console.error("Ошибка при выходе:", err);
                } finally {
                    window.location.href = "/login";
                }
            });
        }
    }

    // Точка входа приложения
    function startApp() {
        initTabs();
        initContainerActions();
        initCreateModal();
        initDeleteModal();
        initGlobalEvents();
        initContainersList();
    }

    startApp();
});
