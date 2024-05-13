export default () => ({
    show: false,
    message: "",
    init: function () {
        document.body.addEventListener("showToaster", (evt) => {
            this.message = evt.detail.message;
            this.show = true;

            setTimeout(() => {
                this.show = false;
            }, 5000);
        });
    },
    close: function () {
        this.show = false;
    },
});
