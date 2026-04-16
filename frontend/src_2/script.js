//-----------------------------------------------------------------------
// INTERFACE INTERACTION
//------------------------------------------------------------------------

// CONTENT HEADER

const contentHeaderScroll = document.querySelector(".content-header");
const homeButton = document.querySelector(".navbutton-home");

// Horizontal scroll for .content-header
contentHeaderScroll.addEventListener("wheel", (e) => {
    e.preventDefault();     // prevent vertical scroll
    contentHeaderScroll.scrollLeft += e.deltaY;  // horizontal scroll
});

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

// loads cities list on sidebar.  
function loadCitiesList() {
    fetch("http://192.168.1.232:8080/cities")
        .then(response => response.json())
        .then(data => {
            const sidebar = document.querySelector(".sidebar-content");

            data.forEach(city => {
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

                sidebar.appendChild(button);
            });
        })
        .catch(error => console.error('Error:', error));
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
            const contentMain = document.querySelector(".content-main");
            contentMain.innerHTML = "";

            // conditions overview swell / ocean info
            const overview = document.createElement("div");
            overview.className = "surf-conditions-overview";

            // data formatting
            let swellHeightft = (parseFloat(data.DomSwellHeightM) * 3.2808).toFixed(1);
            // c to f
            let waterTemp = cToF(data.WaterTempDegC).toFixed(1);
            let airTemp = cToF(data.AirTempDegC).toFixed(1);

            overview.innerHTML = `
                <h3>${spotName} - Current Conditions</h3>
                <p> Swell Height: ${swellHeightft} ft @ ${data.DominantWavePeriodSec} sec</p> 
                <p> Swell Direction: ${data.DomSwellDir}°</p>
                <p> Water Temp: ${waterTemp} °f</p>
                <p> Air Temp: ${airTemp} °f</p>
                <p> Wind: ${data.WindSpeedMph} - (${data.WindDirection})</p>
                <p> Cloud Coverage: ${data.CloudCoverage}</p>
                <p> Chance of Precipitation: ${data.Precipitation}</p>
            `;

            contentMain.appendChild(overview);
        })
        .catch(error => console.error('Error:', error));
};

function cToF(c) {
    return c * 9 / 5 + 32;
};


homeButton.addEventListener("click", () => { resetHomeUI(); });
document.addEventListener("DOMContentLoaded", () => {
    loadCitiesList();
});