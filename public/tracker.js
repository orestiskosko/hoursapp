/**
 * @param {string} startedAt
 */
export default (startedAt) => ({
    /** @type {Date | null} */
    date: startedAt ? new Date(startedAt) : new Date(),

    seconds: 0,

    /** @type {number} */
    interval: null,

    /**
     * Creates the timer display string based on the number of seconds.
     *
     * @returns {string}
     */
    getTimerDisplay() {
        return `${("0" + Math.floor(this.seconds / 3600)).slice(-2)}:${(
            "0" + Math.floor((this.seconds % 3600) / 60)
        ).slice(-2)}:${("0" + (this.seconds % 60)).slice(-2)}`;
    },

    init() {
        console.log("Date", this.date);

        if (this.date) {
            const currentEpoch = new Date().valueOf();
            const secondsDiff = Math.floor(
                (currentEpoch - this.date.valueOf()) / 1000
            );
            this.seconds += secondsDiff;
        } else {
            this.date = new Date();
        }

        this.interval = setInterval(() => {
            this.seconds++;
            document.title = this.getTimerDisplay();
        }, 1000);
    },

    destroy() {
        clearInterval(this.interval);
    },
});
