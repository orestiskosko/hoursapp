/** 
 * @param {string} id
 * @param {boolean} on 
 * @param {string} startedAt 
 */
export default (id, on, startedAt) => ({
    id: id,

    on: on,

    /** @type {Date | null} */
    startedAt: startedAt ? new Date(startedAt) : null,

    seconds: 0,

    /** @type {number} */
    interval: null,

    _getIntervalHandler() {
        return () => {
            this.seconds++
            document.title = this.getTimerDisplay()
        }
    },

    _startTimer() {
        this.interval = setInterval(this._getIntervalHandler(), 1000)
    },

    _stopTimer() {
        clearInterval(this.interval)
    },

    /**
     * Creates the timer display string based on the number of seconds.
     * 
     * @returns {string}
     */
    getTimerDisplay() {
        return `${('0' + Math.floor(this.seconds / 3600)).slice(-2)}:${('0' + Math.floor((this.seconds % 3600) / 60)).slice(-2)}:${('0' + this.seconds % 60).slice(-2)}`
    },

    init() {
        console.log("ID", this.id, "ON", this.on, "StartedAt", this.startedAt)

        if (this.startedAt) {
            const currentEpoch = new Date().valueOf()
            const secondsDiff = this.on ? Math.floor((currentEpoch - this.startedAt.valueOf()) / 1000) : 0
            this.seconds += secondsDiff
        } else {
            this.startedAt = new Date()
        }

        if (this.on) {
            this._startTimer()
        }
    },

    toggle() {
        if (this.id.length === 0) {
            return
        }

        this.on = !this.on
        if (this.on) {
            this._startTimer()
        } else {
            this._stopTimer()
        }
    },

    reset() {
        this.on = false
        this.seconds = 0
        clearInterval(this.interval)
    },

    destroy() {
        clearInterval(this.interval)
    }
})