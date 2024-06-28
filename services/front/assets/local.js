document.addEventListener("DOMContentLoaded", ready);

function ready(){
    standardsTableToggler()
    filterTable()
}

function standardsTableToggler() {
    [].forEach.call( document.querySelectorAll( ".js-table-toggler" ), function ( ele ) {

        ele.addEventListener('click', function(eve){
            var toggleDisplay = eve.target.dataset.toggle;

            [].forEach.call( document.querySelectorAll('.'+toggleDisplay), function(info) {
                if (info.style.display == 'none' || info.style.display == "") {
                    info.style.display = 'table-cell'
                } else {
                    info.style.display = 'none'
                }
                return false
            })

            return false
        }, false)
    } )
}

function filterTable() {
    [].forEach.call(document.querySelectorAll(".js-table-filter"), function(ele){
        ele.addEventListener("change", function(eve){
            var show = eve.target.dataset.show;
            var hide = eve.target.dataset.hide;
            [].forEach.call(document.querySelectorAll(hide), function(i) {
                i.style.display = 'none'
            });
            [].forEach.call(document.querySelectorAll(show), function(i) {
                i.style.display = 'table-row'
            });

        })
    } )

}
