//-----------------------------------------------------------------------
// INTERFACE INTERACTION
//------------------------------------------------------------------------

// ---------------- const variables for script ----------------------
const contentHeaderScroll = document.querySelector(".content-header");
const homeButton = document.querySelector(".navbutton.home-icon");





// hold cities for search filtering
let masterCitiesList = [];

// CONTENT HEADER
// Horizontal scroll for .content-header
contentHeaderScroll.addEventListener("wheel", (e) => {
    e.preventDefault();     // prevent vertical scroll
    contentHeaderScroll.scrollLeft += e.deltaY;  // horizontal scroll
});


citiesSearchInput.addEventListener('input', (e) => {

    let searchValue = e.target.value.toLowerCase().trim();
    console.log('Searching for: ', searchValue);

    if (searchValue === "") {
        renderCitiesList(masterCitiesList);
        return;
    }

    const searchMatch = [];
    masterCitiesList.forEach(city => {
        if (city.name.toLowerCase().includes(searchValue)) {
            searchMatch.push(city);
        };
    });
    renderCitiesList(searchMatch);
});


function fetchCities() {
    fetch("http://192.168.1.232:8080/cities")
        .then(response => response.json())
        .then(data => {
            masterCitiesList = data;
            renderCitiesList(masterCitiesList);
        })
        .catch(error => console.error('Error:', error));
};

function renderCitiesList(citiesList) {

    // clear content before adding cities.
    sidebarContent.innerHTML = "";

    citiesList.forEach(city => {
        const button = document.createElement("button");
        button.className = "city-button";
        button.textContent = `${city.name}, ${city.state.slice(0, 2)}`;
        button.dataset.id = `${city.id}`;


        button.addEventListener('click', function () {
            console.log(".sidebar-content clicked -", button.textContent);

            // clear content-main when new city is selected.
            const contentMain = document.querySelector(".content-main");
            contentMain.innerHTML = "";

            loadSurfSpots(button.dataset.id);
        });

        sidebarContent.appendChild(button);
    });
}


function resetHomeUI() {
    console.log("Home Button clicked.")

    // 1. Clear surf spots list.
    const surfSpotList = document.querySelector(".surf-spot-list");
    surfSpotList.innerHTML = "";

    // 2. Clear main content area.
    const mainContent = document.querySelector(".content-main");
    mainContent.innerHTML = "";

    // 3. Display Go Surf logo.
    const logoDiv = document.createElement("div");
    logoDiv.className = "logodiv";

    const logoText = document.createElement("div");
    logoText.className = "logotext";
    logoText.textContent = "Go Surf";

    logoDiv.appendChild(logoText);
    mainContent.appendChild(logoDiv);
};


// loads surfspots into main content header.
// Parent child relationship from city -> surfspots.
function loadSurfSpots(cityID) {
    console.log("fetching surfspots - cityID - ", cityID);

    fetch(`http://192.168.1.232:8080/surfspots/${cityID}`)
        .then(response => response.json())
        .then(data => {
            const surfSpotList = document.querySelector(".surf-spot-list");
            surfSpotList.innerHTML = "";

            data.forEach(spot => {
                const button = document.createElement("button");
                // setting class name for css manipulation.
                button.className = "spot-button";
                button.textContent = spot.name;
                button.dataset.cityID = spot.cityId;

                button.addEventListener(`click`, function () {

                    console.log("surf-spot-list clicked- ", button.textContent);
                    loadCurrentSurfConditions(spot.id, spot.name);
                });

                surfSpotList.appendChild(button);

                // animate after paint
                requestAnimationFrame(() => {
                    button.classList.add("show");
                });
            });
        });
};

function loadCurrentSurfConditions(spotID, spotName) {
    console.log("fetching current surf conditions -", spotID);

    fetch(`http://192.168.1.232:8080/surfforecast/current/${spotID}`)
        .then(response => response.json())
        .then(data => {

            const swellHeight = metersToFeet(data.DomSwellHeightM).toFixed(1);
            const waterTemp = cToF(data.WaterTempDegC).toFixed(1);
            const airTemp = cToF(data.AirTempDegC).toFixed(1);

            // clear content-main
            const contentMain = document.querySelector(".content-main")
            contentMain.innerHTML = `
                <div class="current-conditions-parent">
                    <div class="conditions-card">
                        <div class="conditions-title">
                            ${spotName} - Current Conditions
                        </div>

                        <div class="conditions-content">
                            <div class="conditions-content-left">
                                <div class="conditions-content-left-title">
                                    Ocean Info
                                </div>
                                <div class="content-left-data">
                                    <p>Dominant swell: ${swellHeight} ft @ ${data.DominantWavePeriodSec} sec</p>
                                    <p>Swell Direction: ${data.DomSwellDir}°</p>
                                    <p>Water Temp: ${waterTemp}°</p>
                                </div>
                            </div>

                            <div class="conditions-content-right">
                                <div class="conditions-content-right-title">
                                    Weather Info
                                </div>
                                <div class="content-right-data">
                                    <p>Air Temp: ${airTemp}°</p>
                                    <p>Wind: ${data.WindSpeedMph} - ${data.WindDirection}</p>
                                    <p>Cloud Coverage: ${data.CloudCoverage}</p>
                                    <p>Precipitation: ${data.Precipitation}%</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            `;

        });
};

function cToF(c) {
    return c * 9 / 5 + 32;
};

function metersToFeet(c) {
    return c * 3.28084
}

homeButton.addEventListener("click", () => { resetHomeUI(); });
document.addEventListener("DOMContentLoaded", () => {
    const sidebarContent = document.querySelector(".sidebar-content");
    const citiesSearchInput = document.querySelector(".location-searchbox input");
    fetchCities();
});