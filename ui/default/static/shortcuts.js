var previousLink, nextLink;

const init = () => {
    previousLink = document.querySelector(".previous");
    nextLink = document.querySelector(".next");
    document.addEventListener("keyup", handleKeyUp);
}

const handleKeyUp = (e) => {
    switch (e.key) {
        case "ArrowLeft":
            if (previousLink) {
                previousLink.click()
            }
            break;
        case "ArrowRight":
            if (nextLink) {
                nextLink.click();
            }
            break;
    }
}

window.addEventListener("load", init);