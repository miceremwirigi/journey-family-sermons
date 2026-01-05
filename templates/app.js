document.addEventListener("click", (e) => {
    const row = e.target.closest('.clickable_row');
    if(row && row.dataset.href) {
        window.location.href = row.dataset.href;
    }
});

let allVideos = [];
let currentPage = 1;
const rowsPerPage = 5;

async function loadVideos() {
    const spinner = document.querySelector(".spinner-container");
    const tableBody = document.querySelector(".sermon-list tbody");

    if (spinner) spinner.style.display = 'flex';

    try {
        const response = await fetch('/videos');
        allVideos = await response.json(); // Store all data globally
        renderTable();
    } catch (error) {
        console.error("Error loading videos: ", error);
        tableBody.innerHTML = "<tr><td colspan='5'> Error loading videos.</td></tr>";
    } finally {
        if (spinner) spinner.style.display = 'none';
    }
}

function renderTable() {
    const tableBody = document.querySelector(".sermon-list tbody");
    tableBody.innerHTML = "";

    // Calculate start and end indices
    const start = (currentPage - 1) * rowsPerPage;
    const end = start + rowsPerPage;
    const paginatedItems = allVideos.slice(start, end);

    paginatedItems.forEach((video, index) => {
        const globalIndex = start + index + 1;
        const row = document.createElement("tr");
        row.className = "clickable_row";
        // row.dataset.href = `https://www.youtube.com/watch?v=${video.ID}`;

        row.innerHTML = `
            <!-- <td>${globalIndex}</td> -->
            <td>
                <iframe width="200" height="113"
                    src="https://www.youtube.com/embed/${video.ID}"
                    frameborder="0" allowfullscreen>
                </iframe>
            </td>
            <td data-label="Title">${video.Title}</td>
            <td data-label="Date">${new Date(video.PublishedAt).toLocaleDateString('en-GB')}</td>
            <td data-label="Action">
                <button class="download__mp3__btn"
                    onclick="event.stopPropagation(); downloadmp3(event, '${video.ID}')">
                    Download mp3
                </button>
            </td>
            <td data-label="Action2">
                <button class="download__mp3__btn"
                    onclick="event.stopPropagation(); downloadmp4(event, '${video.ID}')">
                    Download mp4
                </button>
            </td>
        `;
        tableBody.appendChild(row);
    });

    renderPagination();
}

function renderPagination() {
    const controls = document.getElementById("pagination-controls");
    const totalPages = Math.ceil(allVideos.length / rowsPerPage);
    controls.innerHTML = "";

    // Prev Button
    const prevBtn = document.createElement("button");
    prevBtn.innerText = "«";
    prevBtn.disabled = currentPage === 1;
    prevBtn.onclick = () => { currentPage--; renderTable(); };
    controls.appendChild(prevBtn);

    // Page Numbers
    for (let i = 1; i <= totalPages; i++) {
        const btn = document.createElement("button");
        btn.innerText = i;
        if (i === currentPage) btn.className = "active";
        btn.onclick = () => { currentPage = i; renderTable(); };
        controls.appendChild(btn);
    }

    // Next Button
    const nextBtn = document.createElement("button");
    nextBtn.innerText = "»";
    nextBtn.disabled = currentPage === totalPages;
    nextBtn.onclick = () => { currentPage++; renderTable(); };
    controls.appendChild(nextBtn);
}


const menu = document.querySelector('#mobile-menu');
const menuLinks = document.querySelector('.navbar__menu');

menu.addEventListener('click', function() {
    menu.classList.toggle('is-active');
    menuLinks.classList.toggle('active');
});

function init() {
    loadVideos()
}

if (document.readyState !== "loading"){
    init()
} else {
    document.addEventListener("DOMContentLoaded", init)
}

async function downloadmp3(e, videoId){
    const btn = e.currentTarget;
    btn.disabled = true;
    btn.innerText = "Converting..."

    const fetchlink = `download/mp3/${videoId}`;
    try{
        const data = await fetchWithTimeout(fetchlink, { timeout: 12000 });

        if (!data || data.status !=='tunnel') {
            console.log("Initial fetch url failed or invalid status")
        }

        window.location.href = data.url

        // Artificial delay to keep button disabled while user sees converting ...
        await new Promise(resolve => setTimeout(resolve, 4000))
        
    } catch (error) {
        console.error("Fetch failed:", error.message)
        alert("Error: " + error.message)
    } finally {
        btn.disabled = false;
        btn.innerText = "Convert to mp3"
    }
}

async function fetchWithTimeout(resource, options = {}) {
    const { timeout = 8000 } = options;

    const controller = new AbortController();
    const id = setTimeout(() => controller.abort(), timeout)

    try {
        const response = await fetch(resource, {
            ...options,
        signal: controller.signal  // Connedt the signal
        });
        clearTimeout(id);

        
        // check for status 200-299
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        return await response.json()
    } catch(error) {
        if (error.name ==='AbortError') {
            console.error("Request timed out");
        } else {
            console.error("Request failed:", error);
        }
    }
}

async function downloadmp4(e, videoId){
    const btn = e.currentTarget;
    btn.disabled = true;
    btn.innerText = "Converting..."

    const fetchlink = `download/mp4/${videoId}`;
    try{
        const data = await fetchWithTimeout(fetchlink, { timeout: 12000 });

        if (!data || data.status !=='tunnel') {
            console.log("Initial fetch url failed or invalid status")
        }

        window.location.href = data.url

        // Artificial delay to keep button disabled while user sees converting ...
        await new Promise(resolve => setTimeout(resolve, 4000))
        
    } catch (error) {
        console.error("Fetch failed:", error.message)
        alert("Error: " + error.message)
    } finally {
        btn.disabled = false;
        btn.innerText = "Convert to mp4"
    }
}