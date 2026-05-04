// CONFIG -----------------------------------------------
const API_BASE = "http://localhost:8080";

// Application state -------------------------------------
const state = {
  cities: [],
};

const cloudDict = {
  FEW: "mostly clear",
  SCT: "scattered clouds",
  BKN: "mostly cloudy",
  OVC: "overcast",
  SKC: "clear",
  CLR: "clear",
};

// DOM selectors ------------------------------------------
const DOM = {
  contentHeader: document.querySelector(".content-header"),
  homeButton: document.querySelector(".navbutton.home-icon"),
  sidebarContent: document.querySelector(".sidebar-content"),
  citiesSearchInput: document.querySelector(".location-searchbox"),
  surfSpotList: document.querySelector(".surf-spot-list"),
  contentMain: document.querySelector(".content-main"),
};

// EVENT HANDLERS ------------------------------------

// Horizontal scroll
DOM.contentHeader.addEventListener("wheel", (e) => {
  e.preventDefault();
  DOM.contentHeader.scrollLeft += e.deltaY;
});

// Search filtering
DOM.citiesSearchInput.addEventListener("input", (e) => {
  const searchValue = e.target.value.toLowerCase().trim();

  if (!searchValue) {
    renderCitiesList(state.cities);
    return;
  }

  const filtered = state.cities.filter((city) =>
    city.name.toLowerCase().includes(searchValue),
  );

  renderCitiesList(filtered);
});

// Home button
DOM.homeButton.addEventListener("click", (e) => {
  console.log("Clicked home button.");
  resetHomeUI();
});

// API -------------------------------------------------
function fetchCities() {
  fetch(`${API_BASE}/cities`)
    .then((res) => res.json())
    .then((data) => {
      state.cities = data;
      renderCitiesList(data);
    })
    .catch((err) => console.error("Error fetching cities", err));
}

function fetchSurfSpots(cityId) {
  return fetch(`${API_BASE}/surfspots/${cityId}`).then((res) => res.json());
}

function fetchSurfConditions(spotId) {
  return fetch(`${API_BASE}/surfforecast/current/${spotId}`).then((res) =>
    res.json(),
  );
}

// RENDER -------------------------------------------------
function renderCitiesList(cities) {
  DOM.sidebarContent.innerHTML = "";
  cities.sort();
  cities.forEach((city) => {
    const button = document.createElement("button");
    button.className = "city-button";
    button.textContent = `${city.name}, ${city.state.slice(0, 2)}`;

    button.addEventListener("click", () => {
      DOM.contentMain.innerHTML = "";
      loadSurfSpots(city.id);
    });

    DOM.sidebarContent.appendChild(button);
  });
}

function loadSurfSpots(cityId) {
  fetchSurfSpots(cityId).then((data) => {
    DOM.surfSpotList.innerHTML = "";

    data.forEach((spot) => {
      const button = document.createElement("button");
      button.className = "spot-button";
      button.textContent = spot.name;

      button.addEventListener("click", () => {
        loadCurrentSurfConditions(spot.id, spot.name);
      });

      DOM.surfSpotList.appendChild(button);
      /*
            requestAnimationFrame(() => {
                button.classList.add("show");
            });
            */
    });

    // show first surfspot data
    loadCurrentSurfConditions(data[0].id, data[0].name);
  });
}

function loadCurrentSurfConditions(spotId, spotName) {
  fetchSurfConditions(spotId).then((data) => {
    var swellHeight = display(metersToFeet(data.DomSwellHeightM).toFixed(1));
    var waterTemp = display(cToF(data.WaterTempDegC).toFixed(1));
    var airTemp = display(data.AirTempDegF);
    var windSpeed = display(data.WindSpeedMph);
    var windDir = display(degreesToCardinal(data.WindDirection));
    var cloudCoverage = cloudCoverageDescription(data.CloudCoverage);

    var precipitation = display(data.Precipitation);
    if (precipitation != "NA") {
      precipitation = precipitation + "%";
    }

    DOM.contentMain.innerHTML = `
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
                                <p>Wind: ${windSpeed} mph - ${windDir}</p>
                                <p>Cloud Coverage: ${cloudCoverage}</p>
                                <p>Precipitation: ${precipitation}</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        `;
  });
}

function resetHomeUI() {
  DOM.surfSpotList.innerHTML = "";
  DOM.contentMain.innerHTML = `
        <div class="logodiv">
            <div class="logotext">Go Surf</div>
        </div>
    `;
}

// UTILITY -------------------------------------------------
function cToF(c) {
  return (c * 9) / 5 + 32;
}

function metersToFeet(m) {
  return m * 3.28084;
}

function display(value) {
  if (value === null || value === "") {
    return "NA";
  } else {
    return value;
  }
}

function cloudCoverageDescription(key) {
  return cloudDict[key] || "NA";
}

function degreesToCardinal(deg) {
  if (deg === null) {
    return null;
  } else {
    var directions = ["N", "NE", "E", "SE", "S", "SW", "W", "NW"];
    var index = Math.round((deg / 45) % 8);
    return directions[index];
  }
}

// INIT -------------------------------------------------
fetchCities();
