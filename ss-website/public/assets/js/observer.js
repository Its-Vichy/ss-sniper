const threshold = .30
const options = {
    root: null,
    rootMargin: '0px',
    threshold: threshold
}

const handleIntersect = function (entries, observer) {
    entries.forEach(function (entry) {
        if (entry.intersectionRatio > threshold) {
            entry.target.classList.add('reveal-visible')
            entry.target.classList.remove('reveal')
            observer.unobserve(entry.target)
        }
    })
}

const handleIntersectzoom = function (entries, observer) {
    entries.forEach(function (entry) {
        if (entry.intersectionRatio > threshold) {
            entry.target.classList.add('zoomeffect')
            entry.target.classList.remove('reveal_zoom')
            observer.unobserve(entry.target)
        }
    })
}

document.documentElement.classList.add('reveal-loaded')

window.addEventListener("DOMContentLoaded", function () {
    const observer = new IntersectionObserver(handleIntersect, options)
    const targets = document.querySelectorAll('.reveal')
    targets.forEach(function (target) {
        observer.observe(target)
    });

    const observerzoom = new IntersectionObserver(handleIntersectzoom, options)
    const targets2 = document.querySelectorAll('.reveal_zoom')

    targets2.forEach(function (target) {
        observerzoom.observe(target)
    });
})