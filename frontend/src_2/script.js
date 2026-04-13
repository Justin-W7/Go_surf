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

function fetchSurfSpots(cityID) {
    console.log("fetching surfspots - cityID - ", cityID)
    fetch(`http://localhost:8080/surfspots/${cityID}`)
        .then(response => response.json())
        .then(data => {
            const surfSpotList = document.querySelector(".surf-spot-list");
            surfSpotList.innerHTML = "";

            data.forEach(spot => {
                const button = document.createElement("button");
                button.className = "spot-button";
                button.textContent = `${spot.Name}`;
                button.dataset.cityID = `${spot.cityID}`

                button.addEventListener(`click`, function () {
                    console.log(".surf-spot-list clicked- ", button.textContent);
                });

                surfSpotList.appendChild(button);

                // animate after paint
                requestAnimationFrame(() => {
                    button.classList.add("show");
                });
            });
        });

};

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
                    fetchSurfSpots(button.dataset.id);
                });

                sidebar.appendChild(button);
            });
        })
        .catch(error => console.error('Error:', error));
};

document.addEventListener("DOMContentLoaded", fetchCitiesList);