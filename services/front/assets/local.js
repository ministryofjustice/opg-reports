document.addEventListener("DOMContentLoaded", ready);

function ready(){
    standardsToggler()
}

function standardsToggler() {
    [].forEach.call( document.querySelectorAll( ".toggler" ), function ( ele ) {

        ele.addEventListener('click', function(eve){
            console.log(eve)
            var toggleDisplay = eve.target.dataset.toggle;

            [].forEach.call( document.querySelectorAll('.'+toggleDisplay), function(info) {
                console.log(toggleDisplay)
                console.log(info)
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
