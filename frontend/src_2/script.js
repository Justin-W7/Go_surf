// CONTENT HEADER

// Horizontal scroll for .content-header
const contentHeaderScroll = document.querySelector(".content-header")

contentHeaderScroll.addEventListener("wheel", (e) => {
    e.preventDefault();     // prevent horizontal scroll
    contentHeaderScroll.scrollLeft += e.deltaY;  // horizontal scroll
});
