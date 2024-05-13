// @ts-nocheck
import alpineTracker from "./tracker.js";
import alpineToaster from "./toaster.js";

Date.prototype.addDays = function (days) {
    var date = new Date(this.valueOf());
    date.setDate(date.getDate() + days);
    return date;
};

document.addEventListener("alpine:init", () => {
    Alpine.data("tracker", alpineTracker);
    Alpine.data("toaster", alpineToaster);
});