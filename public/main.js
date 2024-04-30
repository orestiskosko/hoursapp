// @ts-nocheck
import alpineTracker from './track.js'

Date.prototype.addDays = function (days) {
    var date = new Date(this.valueOf());
    date.setDate(date.getDate() + days);
    return date;
}

document.addEventListener('alpine:init', () => {
    Alpine.data('timer', alpineTracker)
})