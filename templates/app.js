// ******************Page rendering logic**********************//
let allVideos = [];
let allPlaylists = []
let currentPage = 1;
const rowsPerPage = 5;
const navBarDefaultItems = `
            <ul class="navbar__menu" id="navbar__menu">
                <li class="navbar__item">
                    <a href="#" id="navlist-toggle" class="navbar__links">
                        Setlists
                    </a>
                </li>
            `;
const debouncedSearch = debounce(() => {
    const val = document.getElementById("video-search").value;
    if (val === "") {
        handleSearch(); // run immediately if input is empty or cleared
    }else {
        handleSearch();
    }
}, 700);

document.addEventListener("click", (e) => {
    const row = e.target.closest('.clickable_row');
    if(row && row.dataset.href) {
        window.location.href = row.dataset.href;
    }
});


// fetches all sermon videos and renders videos list taable
async function loadSermonVideos() {
    const spinner = document.querySelector(".spinner-container");
    const tableBody = document.querySelector(".sermon-list tbody");

    if (spinner) spinner.style.display = 'flex';

    try {
        const response = await fetch('/videos');
        allVideos = await response.json(); // Store all data globally
        if (response.ok) {
            renderTable();
        } else {
            tableBody.innerHTML = `
                    <tr><td colspan='5'> <a href="/admin" style="color: azure;"' >SYNC SERMONS </a> to load videos</td></tr>
                `;
            return
        }
    } catch (error) {
        console.error("Error loading videos: ", error);
    } finally {
        if (spinner) spinner.style.display = 'none';
    }
}

    // Toggle Controller ti switch video table list display from playlist to playlist
function toggleListDisplay() {
    const listToggle = document.getElementById('navlist-toggle');
    const navbarLogo = document.getElementById('navbar__logo');

    listToggle.addEventListener("click", async (e) => {
        e.preventDefault();
        
        // If the button says "Sermons", it means we are in a playlist and want to go BACK to main view
        if (listToggle.innerText === "Sermons") {
            listToggle.innerText = "Setlists";
            navbarLogo.innerText = "Sermons";
            loadSermonVideos(); // Reload all videos
            resetNavbarList();  // Clean up the dynamic playlist items
            // Close mobile menu if open
            const menu = document.querySelector('#mobile-menu');
            if (menu && menu.classList.contains('is-active')) {
                menu.click();
            }
        } else {
            // Otherwise, we are in main view and want to SEE the playlists
            await populatePlaylists(); 
        }
    });
}

// Populate Playlists (with Async Await)
async function populatePlaylists() {
    const navbarMenu = document.getElementById('navbar__menu');
    const listToggle = document.getElementById('navlist-toggle');
    const navbarLogo = document.getElementById('navbar__logo');

    // Wait for playlists to load before trying to loop
    await loadPlaylists(); 
    
    // Clear existing dynamic items to prevent duplicates if clicked twice
    resetNavbarList();

    allPlaylists.forEach((playlist) => {
        const li = document.createElement("li");
        li.className = "navbar__item dynamic-playlist"; 
        li.innerHTML = `<a href="#" class="navbar__links">${playlist.title}</a>`;
        
        li.addEventListener("click", async (e) => {
            e.preventDefault();

            // Close mobile menu if open
            const menu = document.querySelector('#mobile-menu');
            if (menu && menu.classList.contains('is-active')) {
                menu.click();
            }
            
            // UI Feedback
            listToggle.innerText = "Sermons"; // Change toggle text to allow "Back" action
            navbarLogo.innerText = playlist.title; // Show current playlist name as logo
            const tableBody = document.querySelector(".sermon-list tbody"); // clear the existing videos from view
            tableBody.innerHTML = "";
            
            // Load specific items
            await loadPlaylistItems(playlist.id);
        });
        navbarMenu.appendChild(li);
    });
}

// Helper to Reset Menu
function resetNavbarList() {
    // Remove all items with the 'dynamic-playlist' class
    const dynamicItems = document.querySelectorAll('.dynamic-playlist');
    dynamicItems.forEach(item => item.remove());
}

// Awaitable fetch of available playlists 
async function loadPlaylists() {
    // Only fetch if we don't have them or it's an empty array
    if (!allPlaylists || allPlaylists.length === 0) {
        try {
            const response = await fetch("/playlists");
            allPlaylists = await response.json();
        } catch (error) {
            console.error("Error loading playlists:", error);
        }
    }
}

// Fetches all the videos in a playlist of the given id
async function loadPlaylistItems(playlistId) {
    const spinner = document.querySelector(".spinner-container");
    const tableBody = document.querySelector(".sermon-list tbody");

    if (spinner) spinner.style.display = 'flex';

    try {
        const response = await fetch(`/playlists/videos/${playlistId}`);
        allVideos = await response.json(); // Store all data globally
        if (allVideos.length === 0) {
        tableBody.innerHTML = `
                <tr><td colspan='5'> <a href="/admin" style="color: azure;"' >SYNC</a> playlist to view videos.</td></tr>
            `;
        return
        }
        renderTable();
    } catch (error) {
        console.error("Error loading videos: ", error);
        tableBody.innerHTML = "<tr><td colspan='5'> Error loading videos.</td></tr>";
    } finally {
        if (spinner) spinner.style.display = 'none';
    }
}

// function addsetlistitem() {
//     const navbarMenu = document.getElementById("navbar__manu")
//     loadedMenus
// }

function renderTable(dataToRender = allVideos) {
    const tableBody = document.querySelector(".sermon-list tbody");
    if (!tableBody) return;
    tableBody.innerHTML = "";

    // Calculate start and end indices
    const start = (currentPage - 1) * rowsPerPage;
    const end = start + rowsPerPage;
    const paginatedItems = dataToRender.slice(start, end);

    if (paginatedItems.length === 0) {
        const isSearching = document.getElementById("video-search").value !== "";

        tableBody.innerHTML = isSearching
            ? "<tr><td colspan='5' style='text-align:center; padding: 20px;'>🔍 videos matching your search criteria</td></tr>"
            : "<tr><td colspan='5' style='text-align:center; padding: 20px;'>No videos available in this list</td></tr>";
            
        const controls = DocumentTimeline.getElementById("pagination-controls"); // clear pagination if no item is found
        if (controls) controls.innerHTML="";
        return
    }
    paginatedItems.forEach((video, index) => {
        const globalIndex = start + index + 1;
        const row = document.createElement("tr");
        row.className = "clickable_row";

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

    renderPagination(dataToRender.length, dataToRender);
}

// ***********************************************************************//


// ***************Pagination Logic*******************************//

function renderPagination(totalItems, dataReference = allVideos) {
    const controls = document.getElementById("pagination-controls");
    const totalPages = Math.ceil(totalItems / rowsPerPage);
    if (totalPages <=1) {
        controls.innerHTML="";
        return
    }

    controls.innerHTML = "";

    const createBtn = (page, text, isDisabled = false, active = false) => {
        const btn = document.createElement("button");
        btn.innerHTML = text;
        btn.disabled = isDisabled;
        if (active) btn.className = "active";
        btn.onclick = () => {
            currentPage = page;
            renderTable(dataReference);
            window.scrollTo({ top: 0, behavior: 'smooth'});
        }
        return btn;
    };

    controls.appendChild(createBtn(currentPage -1, "«", currentPage === 1));

    const range = []
    const delta = 1; // number of heighbours to show, e.g., 1 showing 2,[3], 4

    for (let i = 1; i <=totalPages; i++) {
        // always show the first, last, and current page + its neighnors
        if (i === 1 || i === totalPages || (i >= currentPage - delta && i <= currentPage + delta)) {
            range.push(i);
        }
    }

    let lastItem;
    range.forEach(i => {
        if (lastItem) {
            if (i - lastItem === 2) {
                // if theres exactly a gap of 1 page, just show that page number
                controls.appendChild(createBtn(lastItem + 1, lastItem + 1))
            } else if (i - lastItem > 2) {
                const dots = document.createElement("span");
                dots.innerHTML = "...";
                dots.className = "pagination-dots";
                controls.appendChild(dots);
            }
        }
        controls.appendChild(createBtn(i, i, false, i === currentPage));
        lastItem = i;
    });

    controls.appendChild(createBtn(currentPage + 1, "»", currentPage === totalPages))

    // Add go to page functionality
    const goContainer = document.createElement("div");
    goContainer.className = "pagination-go";
    goContainer.innerHTML = `
            <input type="number" id="jump-page" min="1" max="${totalPages}" placeholder="Page...">
            <button id="jump-btn">Go</button>
        `;
    controls.appendChild(goContainer);

    document.getElementById("jump-btn").onclick = () => {
        const val = parseInt(document.getElementById("jump-page").value);
        if (val >= 1 && val <= totalPages) {
            currentPage = val;
            renderTable();
            window.scrollTo({ top: 0, behaviour: 'smooth' });
        }
    };
}

const menu = document.querySelector('#mobile-menu');
const menuLinks = document.querySelector('.navbar__menu');

menu.addEventListener('click', function() {
    menu.classList.toggle('is-active');
    menuLinks.classList.toggle('active');
});

function init() {
    if (!window.location.pathname.includes("admin")){
    loadSermonVideos();
    toggleListDisplay();
    document.getElementById("video-search").addEventListener("keypress", (e) => {
        if (e.key === "Enter") {
            e.preventDefault(); // prevent page reload
            handleSearch();
        }
    });
    }
}

if (document.readyState !== "loading"){
    init()
} else {
    document.addEventListener("DOMContentLoaded", init)
}

function updateNavLinks() {
    const path = window.location.pathname;
    
    // Check if we are on the Home page 
    const isHome = path === "/";
    // Check if we are on the Admin page
    const isAdmin = path === "/admin";

    const adminLink = document.querySelector('.admin-nav-item');
    const homeLink = document.querySelector('.home-nav-item'); // Add this class to your Home <li>

    if (isHome) {
        if (adminLink) adminLink.style.display = 'flex';
        if (homeLink) homeLink.style.display = 'none'; // Hide "Home" link if already home
    } else if (isAdmin) {
        if (adminLink) adminLink.style.display = 'none'; // Hide "Admin" link if already in admin
        if (homeLink) homeLink.style.display = 'flex';
    }
}

// ***********************************************************************//

// *************************Admin Specific logic***************************//

// Handle the "Fetch/Update" button triggers
async function triggerUpdate(endpoint) {
    const spinner = document.getElementById("loading-spinner");
    if (spinner) spinner.style.display = 'flex';

    try {
        const response = await fetch(endpoint, {
            method: 'POST', 
            headers: { 'Content-Type': 'application/json' }
        });
        
        if (response.ok) {
            alert("Sync triggered successfully!");
            loadPlaylists();
        } else {
            throw new Error("Update failed");
        }
    } catch (error) {
        alert("Error: " + error.message);
    } finally {
        if (spinner) spinner.style.display = 'none';
    }
}

// Handle adding a new Playlist URL
const playlistForm = document.getElementById('playlist-form');
if (playlistForm) {
    playlistForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const urlInput = document.getElementById('playlist-url-1');
        const url = urlInput.value;

        if (!url) return alert("Please enter a URL");

        try {
            const response = await fetch('/playlists/add', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ url: url })
            });

            if (response.ok) {
                alert("Playlist added for tracking!");
                urlInput.value = ""; // Clear input
                loadPlaylists(); // Refresh list
            } else {
                alert("Failed to add playlist.");
            }
        } catch (error) {
            console.error("Error:", error);
        }
    });
}

// Update the existing loadPlaylists to render into the Admin Table
async function loadPlaylistsAdmin() {
    const tableBody = document.querySelector(".sermon-list tbody");
    const spinner = document.getElementById("loading-spinner");
    
    if (spinner) spinner.style.display = 'flex';

    try {
        const response = await fetch("/playlists");
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const playlists = await response.json();
        if (!playlists || !Array.isArray(playlists)) {
            tableBody.innerHTML = "<tr><td colspan='4'>No playlists found in database.</td></tr>";
            return;
        }
        
        tableBody.innerHTML = ""; 
        playlists.forEach(pl => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td><img src="${pl.thumbnail_url || 'static/favicon.png'}" width="40" style="border-radius:4px"></td>
                <td data-label="Title"><strong>${pl.title}</strong></td>
                <td data-label="Items">${pl.item_count} items</td>
                <td data-label="Actions">
                    <div style="display: flex; gap: 8px;">
                        <button class="download__mp3__btn" 
                                onclick="event.stopPropagation(); syncPlaylist('${pl.id}', '${pl.title}')">
                            Sync Items
                        </button>
                        
                        <button class="download__mp3__btn" 
                                style="background-color: #8b0000;" 
                                onclick="event.stopPropagation(); deletePlaylist('${pl.id}')">
                            Delete
                        </button>
                    </div>
                </td>
            `;
            tableBody.appendChild(row);
        });
    } catch (error) {
        console.error("Fetch error:", error);
        tableBody.innerHTML = "<tr><td colspan='4'>Error loading playlists.</td></tr>";
    } finally {
        if (spinner) spinner.style.display = 'none';
    }
}

async function deletePlaylist(playlistId) {
    if (!confirm("Are you sure you want to stop tracking this playlist and remove it from the repository?")) {
        return;
    }

    const spinner = document.getElementById("loading-spinner");
    if (spinner) spinner.style.display = 'flex';

    try {
        const response = await fetch(`/playlists/${playlistId}`, {
            method: 'DELETE', // Adjust to 'POST' if your backend doesn't support DELETE
            headers: {
                'Content-Type': 'application/json'
            }
        });

        if (response.ok) {
            alert("Playlist removed successfully.");
            loadPlaylistsAdmin(); // Refresh the table
        } else {
            const errorData = await response.json();
            throw new Error(errorData.message || "Failed to delete");
        }
    } catch (error) {
        console.error("Delete error:", error);
        alert("Error deleting playlist: " + error.message);
    } finally {
        if (spinner) spinner.style.display = 'none';
    }
}

async function syncPlaylist(playlistId) {
    const spinner = document.getElementById("loading-spinner");
    if (spinner) spinner.style.display = 'flex';
    let p = document.createElement('p')
    spinner.appendChild(p)

    try {
        // Calling the endpoint that runs the SyncSinglePlaylist logic in Go
        const response = await fetch(`/playlists/sync/${playlistId}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' }
        });

        if (response.ok) {
            allVideos = [] // clear local cache
            allPlaylists = [] // clear local cache
            alert("Sync successfull! Data will refresh.");
            loadPlaylistsAdmin(); // Refresh to show updated item count
        } else {
            throw new Error("Failed to sync items");
        }
    } catch (error) {
        alert("Sync Error: " + error.message);
    } finally {
        if (spinner) spinner.style.display = 'none';
    }
}

// Initial check: if we are on admin.html, load playlists instead of videos
if (window.location.pathname.includes('admin')) {
    console.log("on admin page")
    document.addEventListener("DOMContentLoaded", loadPlaylistsAdmin());
}

// ***********************************************************************//

// **************************Download logic******************************//

async function downloadmp3(e, videoId){
    const btn = e.currentTarget;
    btn.disabled = true;
    btn.innerText = "Converting..."

    const fetchlink = `download/mp3/${videoId}`;
    try{
        const data = await fetchWithTimeout(fetchlink, { timeout: 30000 });

        if(!data) {
            throw new Error("Server did not respond in time.")
        }

        if (!data || data.status !=='tunnel') {
            console.log("Initial fetch url failed or invalid status");
            throw new Error("Invalid server response status.")
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
// **********************************************************************************************//

// **************************Search and Filter logic****************************************//
function filterVideosByPlaylist(playlistId) {
    const filtered = allVideos.filter(v => v.playlist_id === playlistId);
    renderTable(filtered);
}

function handleSearch() {
    const searchInput = document.getElementById("video-search");
    if (!searchInput) return;

    const searchTerm = document.getElementById("video-search").value.toLowerCase();

    if (!allVideos || !Array.isArray(allVideos) || allVideos.length ===0) {
        console.warn("Search attempted before videos were loaded")
        return;
    }

    const filteredVideos = allVideos.filter(video => {
        const title = (video.Title || "").toLowerCase();
        const description = (video.Description || "").toLowerCase();

        const matchesTilte = title.includes(searchTerm);

        const matchesDescription = description.includes(searchTerm);

        return matchesTilte || matchesDescription;
        });

        currentPage = 1;
        renderTable(filteredVideos);
        console.log(`Search for "${searchTerm}" completed. Found ${filteredVideos.length} results. `);
}

function handleSort() {
    const sortType = document.getElementById("sort-filter").value;

    if (!allVideos || allVideos.length === 0) return;

    allVideos.sort((a, b) => {
        // Ensure we have strings to work with
        const titleA = String(a.Title || "");
        const titleB = String(b.Title || "");

        if (sortType === "newest") {
            return new Date(b.PublishedAt) - new Date(a.PublishedAt);
        } else if (sortType === "oldest") {
            return new Date(a.PublishedAt) - new Date(b.PublishedAt);
        } else if (sortType === "title-asc") {
            const cleanA = titleA.replace(/^\W+/, "").trim(); // remove emojis with regex and any accidental space
            const cleanB = titleB.replace(/^\W+/, "").trim(); // remove emojis with regex
            return cleanA.localeCompare(cleanB, undefined, {
                sensitivity: 'base', // treat a and A the same
                numeric: true // ensure 'Part 2' comes before 'Part 10'
            });
        } else if (sortType === "title-dsc") {
            const cleanA = titleA.replace(/^\W+/, "").trim();     
            const cleanB = titleB.replace(/^\W+/, "").trim();   
            return cleanB.localeCompare(cleanA, undefined, {
                sensitivity: 'base',    
                numeric: true   
            });
        }
    });
    currentPage = 1;
    renderTable();
}

/**
 * @param {Function} func - The function to delay
 * @param {number} delay - TIme in milliseconds to wait
 */
function debounce(func, delay = 300) {
    let timeoutId;
    return (...args) => {
        if (timeoutId) clearTimeout(timeoutId);

        timeoutId = setTimeout(() => {
            func.apply(null, args);
        }, delay);
    };
}

// **********************************************************************************************//
