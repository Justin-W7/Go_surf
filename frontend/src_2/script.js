//-----------------------------------------------------------------------
// INTERFACE INTERACTION
//------------------------------------------------------------------------

// CONTENT HEADER

// Horizontal scroll for .content-header
const contentHeaderScroll = document.querySelector(".content-header")

contentHeaderScroll.addEventListener("wheel", (e) => {
    e.preventDefault();     // prevent horizontal scroll
    contentHeaderScroll.scrollLeft += e.deltaY;  // horizontal scroll
});

// loads cities list on sidebar.  
function fetchCitiesList() {
    fetch("http://localhost:8080/cities")
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

                    fetchSurfSpots(button.dataset.id);
                });

                sidebar.appendChild(button);
            });
        })
        .catch(error => console.error('Error:', error));
};

// loads surfspots into main content header.
// Parent child relationship from city -> surfspots.
function fetchSurfSpots(cityID) {
    console.log("fetching surfspots - cityID - ", cityID);

    fetch(`http://localhost:8080/surfspots/${cityID}`)
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
                    console.log(".surf-spot-list clicked- ", button.textContent);
                    fetchCurrentSurfConditions(spot.id, spot.name);
                });

                surfSpotList.appendChild(button);

                // animate after paint
                requestAnimationFrame(() => {
                    button.classList.add("show");
                });
            });
        });
};

function fetchCurrentSurfConditions(spotID, spotName) {
    console.log("fetching current surf conditions -", spotID);

    fetch(`http://localhost:8080/surfforecast/current/${spotID}`)
        .then(response => response.json())
        .then(data => {
            const contentMain = document.querySelector(".content-main");
            contentMain.innerHTML = "";

            // conditions overview
            const overview = document.createElement("div");
            overview.className = ".surf-conditions-overview";

            overview.innerHTML = `
                <h3>${spotName}</p>
                <h4>${data.DomSwellHeightM}m @ ${data.DomSwellDir} deg</h2>
            `;

            contentMain.appendChild(overview);
        })
        .catch(error => console.error('Error:', error));
};


document.addEventListener("DOMContentLoaded", fetchCitiesList);