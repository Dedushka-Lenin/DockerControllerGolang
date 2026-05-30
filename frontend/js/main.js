// Инициализация интерфейса после загрузки DOM
document.addEventListener("DOMContentLoaded", () => {
    const containersListEl = document.getElementById("containers-list");
    const controlPanelEl = document.getElementById("control-panel");
    const panelTitleEl = document.getElementById("selected-container-title");
    
    // Элементы управления для дальнейшей привязки логики
    const btnCreate = document.getElementById("btn-create");
    const btnLogout = document.getElementById("btn-logout");
    const btnContainer = document.getElementById("btn-container"); // Изменено с btnInfo
    const btnLogs = document.getElementById("btn-logs");           // Новая кнопка
    const btnConsole = document.getElementById("btn-console");

    let selectedContainerId = null;

    // 1. Получение и отрисовка списка контейнеров
    async function initContainers() {
        try {
            const res = await loadData("GET", "/containers/get");
            
            containersListEl.innerHTML = ""; // Очищаем индикатор загрузки

            if (!res.containers_list || res.containers_list.length === 0) {
                containersListEl.innerHTML = '<li class="loading">Контейнеры не найдены</li>';
                return;
            }

            // Рендерим каждый элемент из массива []domain.Container
            res.containers_list.forEach(container => {
                const li = document.createElement("li");
                li.className = "container-item";
                li.textContent = container.Name;
                li.dataset.id = container.Id; // Сохраняем ID в data-атрибут

                // Обработчик клика на контейнер
                li.addEventListener("click", () => {
                    // Подсвечиваем активный элемент
                    document.querySelectorAll(".container-item").forEach(el => el.classList.remove("active"));
                    li.classList.add("active");

                    // Запоминаем ID и открываем правую панель управления
                    selectedContainerId = container.Id;
                    panelTitleEl.textContent = `Управление: ${container.Name}`;
                    controlPanelEl.classList.remove("hidden");
                });

                containersListEl.appendChild(li);
            });

        } catch (err) {
            containersListEl.innerHTML = `<li class="loading" style="color:red">Ошибка: ${err.message}</li>`;
        }
    }

    // Элементы модального окна
    const modalCreate = document.getElementById("modal-create");
    const btnModalClose = document.getElementById("btn-modal-close");
    const btnModalSubmit = document.getElementById("btn-modal-submit");
    const inputName = document.getElementById("input-name"); // Получаем поле имени
    const inputUrl = document.getElementById("input-url");
    const modalError = document.getElementById("modal-error");

    // Открытие модального окна
    btnCreate.addEventListener("click", () => {
        // Очищаем прошлые данные и ошибки при открытии
        inputName.value = "";
        inputUrl.value = "";
        modalError.textContent = "";
        modalError.classList.add("hidden");
        
        modalCreate.classList.remove("hidden");
    });

    // Закрытие модального окна по крестику
    btnModalClose.addEventListener("click", () => {
        modalCreate.classList.add("hidden");
    });

    // Закрытие модального окна при клике на темную область вокруг
    modalCreate.addEventListener("click", (event) => {
        if (event.target === modalCreate) {
            modalCreate.classList.add("hidden");
        }
    });

    // Отправка формы создания
    btnModalSubmit.addEventListener("click", async () => {
        const nameValue = inputName.value.trim();
        const urlValue = inputUrl.value.trim();

        // Валидация перед отправкой
        if (!nameValue) {
            showModalError("Поле Название не может быть пустым");
            return;
        }
        if (!urlValue) {
            showModalError("Поле URL не может быть пустым");
            return;
        }

        try {
            btnModalSubmit.disabled = true; // Блокируем кнопку на время запроса
            modalError.classList.add("hidden");

            // Отправка запроса с обоими параметрами
            await loadData("POST", "/containers/create", { 
                name: nameValue,
                url: urlValue 
            });

            // Если успешно: закрываем окно и обновляем список контейнеров
            modalCreate.classList.add("hidden");
            initContainers(); 
            
        } catch (err) {
            showModalError(err.message || "Не удалось создать контейнер");
        } finally {
            btnModalSubmit.disabled = false;
        }
    });

    // Функция для удобного вывода ошибок
    function showModalError(message) {
        modalError.textContent = message;
        modalError.classList.remove("hidden");
    }

    // 2. Обработчики действий для кнопок управления
    btnContainer.addEventListener("click", () => {
        if (!selectedContainerId) return;
        alert(`Запрос информации о контейнере ID: ${selectedContainerId}`);
        // Вызов: loadData("GET", `/containers/info/${selectedContainerId}`)
    });

    btnLogs.addEventListener("click", () => {
        if (!selectedContainerId) return;
        alert(`Получение логов для контейнера ID: ${selectedContainerId}`);
        // Вызов: loadData("GET", `/containers/logs/${selectedContainerId}`)
    });

    btnConsole.addEventListener("click", () => {
        if (!selectedContainerId) return;
        alert(`Открытие Console для контейнера ID: ${selectedContainerId}`);
    });


    btnLogout.addEventListener("click", async () => {
        try {
            await loadData("POST", "/logout"); 
        } catch (error) {
            console.error("Ошибка при выходе:", error);
        }
        window.location.href = "/login";
    });

    // Запускаем чтение данных при старте страницы
    initContainers();
});
