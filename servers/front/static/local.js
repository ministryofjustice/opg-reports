document.addEventListener("DOMContentLoaded", ready);

function ready(){
    tableToggler()
    tableChoiceFilter()
    tableSorter()
    generateFormFilterOptions()
    tableOptionFilters()
    compareGraphClear()
    comparisonToggler()
    compareChartItemListeners()
    sparklines()
}

// ==========================
// SPARKLINE
// ==========================
/**
 * Create a constructor for sparklines that takes some sensible defaults
 * and merges in the individual chart options.
 */
Highcharts.SparkLine = function (a, b, c) {
    const hasRenderToArg = typeof a === 'string' || a.nodeName;
    let options = arguments[hasRenderToArg ? 1 : 0];
    const defaultOptions = {
        chart: {
            renderTo: (
                (options.chart && options.chart.renderTo) ||
                (hasRenderToArg && a)
            ),
            backgroundColor: null,
            borderWidth: 0,
            type: 'area',
            margin: [2, 0, 2, 0],
            height: 20,
            style: {
                overflow: 'visible'
            },
            // small optimalization, saves 1-2 ms each sparkline
            skipClone: true
        },
        exporting: {enabled: false},
        title: {
            text: ''
        },
        credits: {
            enabled: false
        },
        xAxis: {
            labels: { enabled: false },
            title: { text: null },
            startOnTick: false,
            endOnTick: false,
            tickPositions: []
        },
        yAxis: {
            endOnTick: false,
            startOnTick: false,
            labels: { enabled: false },
            title: { text: null },
            tickPositions: [0]
        },
        legend: {
            enabled: false
        },
        tooltip: {
            enabled: false,
        },
        plotOptions: {
            series: {
                animation: false,
                lineWidth: 1,
                shadow: false,
                states: {
                    hover: {
                        enabled: false,
                    }
                },
                marker: {
                    enabled: false,
                },
                fillOpacity: 0.25
            },
            column: {
                negativeColor: '#910000',
                borderColor: 'silver'
            }
        }
    };

    options = Highcharts.merge(defaultOptions, options);

    return hasRenderToArg ?
        new Highcharts.Chart(a, options, c) :
        new Highcharts.Chart(options, b);
};


function sparklines() {
    let limit = 20
    let delay = 1500
    let intervalT = setInterval(function(){
        done = drawSparklines(limit)
        if (done < limit) {
            clearInterval(intervalT)
        }
    }, delay)
}

function drawSparklines(limit) {
    let sparkSelector = ".js-sparkline:not(.js-sparked)"

    let sparks = [...document.querySelectorAll(sparkSelector)].slice(0, limit)
    sparks.forEach(spark => {
        let row = spark.closest("tr")
        let data = []
        // get the data in the row
        row.querySelectorAll(".data-cell span").forEach( cell => {
            let txt = cell.getAttribute("title").trim()
            let fl = parseFloat(txt)
            data.push( parseFloat ( fl.toFixed(2) ))
        })
        // remove the loaders
        spark.querySelectorAll(".loader").forEach(l => { spark.removeChild(l) })
        Highcharts.SparkLine(spark, {
            series: [{ data: data, pointStart: 1 }],
            chart: {}
        });

        spark.classList.add("js-sparked")
        spark.classList.remove("js-spark-loading")
    })
    return sparks.length
}

// ==========================
// COMPARISON GRAPHS
// ==========================

// compareSetup is run when the charts are shown / hidden
function compareSetup() {
    compareTableDefaultSelection()
    compareGraphInsert()
}

// comparisonToggler adds listeners to hide / show the charts
// - run on page load only
function comparisonToggler() {
    let selector = ".js-compare-enable .js-compare-toggler"
    let itemsSelector = ".js-compare-chart, .js-compare-item, .js-compare-intro"
    let toggle = "js-compare-enabled"
    document.querySelectorAll(selector).forEach(check => {
        check.addEventListener("click", function(event){
            document.querySelectorAll(itemsSelector).forEach(i => { i.classList.toggle(toggle) })
            // - reset charts
            compareGraphClear()
            compareSetup()
        })
    })
}

// compareChartRender sets up the config for the highchart and
// then triggers the draw of that
// - uses table structure (thead th, tbody th) to control the
//      xAxis and series data
function compareChartRender(container, tableId) {
    let headerSelector = "thead .data-cell"
    let itemSelector = "tbody .js-compare-active"
    let seriesSelector = ".data-cell span"
    // stock config, we know get the rest from the table
    var config = {
        title: {text: ""},
        exporting: {enabled: false},
        tooltip: { valuePrefix: "$" }
    }

    // get the x-axis from the dates in the thead
    var xAxis = [];
    document.querySelector("#"+tableId).querySelectorAll(headerSelector).forEach(cell => {
        xAxis.push(cell.textContent)
    })
    // set the xAxis config
    config["xAxis"] = {
        categories: xAxis,
        crosshair: true,
        accessibility: {
            description: "Months"
        }
    }
    // now generate the series data
    var series = []
    // find the activated rows
    document.querySelector("#"+tableId).querySelectorAll(itemSelector).forEach(row => {
        // the name is made from the col header cells
        var name = ""
        var data = []
        row.querySelectorAll("th").forEach(th => { name = name + " "+ th.textContent})
        // the data series is then made from all the data cells
        row.querySelectorAll(seriesSelector).forEach( cell => {
            let txt = cell.getAttribute("title").trim()
            let fl = parseFloat(txt)
            data.push( parseFloat ( fl.toFixed(2) ))
        })
        series.push({ name: name.trim(), data: data})
    })

    config["series"] = series
    // set the seperators
    Highcharts.setOptions({
        lang: {
          decimalPoint: '.',
          thousandsSep: ','
        }
    });
    Highcharts.chart(container, config)

}

// compareGraphClear removes any existing chart and deselectes
// any check boxes
// - run on page load to setup and then everytime the charts are hidden / shown
function compareGraphClear() {
    let sel = ".js-compare-chart"
    let graphs = ".js-compare-graph"
    let checks = ".js-compare-item"
    document.querySelectorAll(sel).forEach(chart => {
        chart.querySelectorAll(graphs).forEach(g => { chart.removeChild(g) })
    })
    document.querySelectorAll(checks).forEach(ch => { ch.checked = false})
}

// compareTableDefaults clisk the first X items in the table for rendering
function compareTableDefaultSelection() {
    let defaultItems = 5
    let charts = ".js-compare-chart"
    let checks = ".js-compare-item"
    document.querySelectorAll(charts).forEach(chart => {
        let dataTableSelector = chart.dataset.compare
        document.querySelectorAll(dataTableSelector).forEach(tbl => {
            tbl.querySelectorAll(checks).forEach((ch, i) => {
                if (i < defaultItems) {
                    ch.checked = true
                    compareTableEventToggleClass(ch)
                }
            })
        })
    })
}
// compareGraphInsert inserts a new container into each chart block
// and then calls the chart to be rendered
function compareGraphInsert() {
    let charts = ".js-compare-chart"
    document.querySelectorAll(charts).forEach(chart => {
        let ts = Date.now()
        let dataTableSelector = chart.dataset.compare
        let dataTable = document.querySelector(dataTableSelector)
        let container = `container-${ts}`
        // now insert one
        var gph = document.createElement("figure")
        gph.className = "js-compare-graph highcharts-figure"
        gph.innerHTML = `<div id="${container}"></div>`
        chart.insertBefore(gph, chart.firstChild)
        // call the renderer
        compareChartRender(container, dataTable.id)
    })
}
// compareTableEventToggleClass used to toggle active class on the parent row
function compareTableEventToggleClass(ch) {
    let activeClass = "js-compare-active"
    ch.closest("tr").classList.toggle(activeClass)

}
// compareChartItemListeners
// - adds listener to checkbox that will toggle a class on each checkbox parent row
// - adds listner to checkbox that then renders the chart (via compareGraphInsert)
// - runs on page load only
function compareChartItemListeners() {
    let chartSelector = ".js-compare-chart"
    let checkboxSelector = ".js-compare-item"

    document.querySelectorAll(chartSelector).forEach(chart => {
        // find the data table
        let now = Date.now();
        let dataTableSelector = chart.dataset.compare
        let dataTables = document.querySelectorAll(dataTableSelector)
        // for each table, we now trigger the first 5 items
        dataTables.forEach(tbl => {
            tbl.setAttribute("id", "js-compare-chart-tbl-"+now)
            let checks = tbl.querySelectorAll(checkboxSelector)
            checks.forEach((ch, i) => {
                let clickTriggerClass = function(event){ compareTableEventToggleClass(ch) }
                let clickTriggerGraph = function(event){ compareGraphClear(); compareGraphInsert(); }
                // remove listeners
                ch.removeEventListener("click", clickTriggerClass)
                ch.removeEventListener("click", clickTriggerGraph)
                // toggle a class
                ch.addEventListener("click",  clickTriggerClass)
                // toggle chart
                ch.addEventListener("click", clickTriggerGraph)
            })
        })
    })
}


// ==========================
// TABLE ITEM VISIBILITY
// ==========================
// Show / hide elements within table rows when + icons clicked
function tableToggler() {
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


// ==========================
// TABLE FILTERS
// ==========================
// Filter table based on radio button selections
function tableChoiceFilter() {
    let choices = document.querySelectorAll(".js-table-choice-filter")
    choices.forEach(ele => {
        ele.addEventListener("change", function(eve){
            let allQ = eve.target.dataset.all;
            let showQ = eve.target.dataset.show;
            let selector = allQ + '['+showQ+']'

            if (showQ != undefined){
                document.querySelectorAll(allQ).forEach(r => {r.style.display = 'none'})
                document.querySelectorAll(selector).forEach(r => {r.style.display = 'table-row'})
            } else {
                document.querySelectorAll(allQ).forEach(r => {r.style.display = 'table-row'})
            }

        })
    } )
}

function tableFilterByAll(tableData, filters){

    tableData.forEach(row => {
        let show = true;
        let cols = row.querySelectorAll('td,th')
        // console.log(cols);
        filters.forEach(filter => {
            let colNum = filter.dataset.col;
            let sel = filter.querySelector('select');
            let val = sel.value;
            let comp = cols[colNum].textContent;
            if (val != "all" && comp != val) {
                show = false;
            }
        })
        if (show == true) {
            row.style.display = 'table-row';
        } else {
            row.style.display = 'none';
        }
    })
}

function tableOptionFilters() {
    let formSelector = `.js-table-filter-options`
    let filterSelector = `.js-table-filter-select`
    let forms = document.querySelectorAll(formSelector)
    forms.forEach(f => {
        let done = f.dataset.filterActive;
        let tableData = document.querySelectorAll(f.dataset.filterrows);
        if (done !== true){
            let filters = f.querySelectorAll(filterSelector)

            filters.forEach(filter => {
                let sel = filter.querySelector('select');
                sel.addEventListener("change", e => tableFilterByAll(tableData, filters))
            })
        }

    })

}
// generateFormFilterOptions adds all the unique values from a column into a select within a form
// allowing dynamic data sets pulled from the page content
function generateFormFilterOptions() {
    let formSelector = `.js-table-filter-options-generate`
    let filterSelect = `.js-table-generate-options`
    let generators = document.querySelectorAll(formSelector)

    generators.forEach(block => {
        let done = block.dataset.generated
        if (done !== true) {
            let dataSource = block.dataset.optionssource;
            let filterWrappers = block.querySelectorAll(filterSelect)
            filterWrappers.forEach(wrapper => {
                let col = wrapper.dataset.col;
                let sel = wrapper.querySelector('select');
                let selectedCols = document.querySelectorAll(dataSource + `[data-col="${col}"]`);
                let colValues = [];
                selectedCols.forEach(c => { colValues.push(c.textContent) })
                let unique = [...new Set(colValues)];
                unique = unique.sort()
                // add the data
                unique.forEach(val => {
                    sel.add(new Option(val, val))
                })
                wrapper.classList.add(`js-table-filter-select`)
            })
            block.dataset.generated = true

        }
    })

}


// ==========================
// TABLE SORTING
// ==========================
// Compare values - strip $ and ,
// - try sorting as a number via float parsing
function tableSortCompareValues(a, b) {
    aVal = a.replace("$", "").replace(",", "");
    bVal = b.replace("$", "").replace(",", "");
    comp = 0
    if (parseFloat(aVal)){
        aVal = parseFloat(aVal);
    }
    if (parseFloat(bVal)){
        bVal = parseFloat(bVal);
    }

    if (aVal < bVal) {
        comp = -1
    } else if (aVal > bVal) {
        comp = 1
    }
    return comp

}
// Handle the sort trigger event and reorder the table body
function tableSortHandler(thead, tbody, th, colNum) {
    let defaultDir = "asc";
    let rows = Array.from(tbody.querySelectorAll(`tr`));
    let sortDir = th.dataset.sortdir;
    if (sortDir == undefined){
        sortDir = defaultDir
    }
    if (sortDir == "asc") {
        sortDir = "desc"
    } else {
        sortDir = "asc"
    }
    // and then just... sort the rows:
    rows.sort( (r1,r2) => {
        // get each row's relevant column
        var t1 = r1.querySelectorAll('td,th')[colNum];
        var t2 = r2.querySelectorAll('td,th')[colNum];
        var c1 = t1.textContent
        var c2 = t2.textContent
        comp = tableSortCompareValues(c1, c2)
        return comp
    });
    if (sortDir == "desc") {
        rows = rows.reverse()
    } else {
        tableSortDesc = true
    }

    thead.querySelectorAll(`th`).forEach(t => {t.classList.remove('sorted')})
    th.classList.add('sorted');
    th.dataset.sortdir = sortDir;
    // and then the magic part that makes the sorting appear on-page:
    rows.forEach(row => tbody.appendChild(row));
}
// Main table sorter call
function tableSorter() {
    tables = document.querySelectorAll("table.js-table-sorter")
    tables.forEach( table => {
        thead = table.querySelector(`thead`);
        tbody = table.querySelector(`tbody`);

        thead.querySelectorAll(`th`).forEach((th, position) => {
            th.addEventListener(`click`, evt => tableSortHandler(thead, tbody, th, position));
        });
        // trigger a click on the last column
        var sel = ".govuk-table__head th:last-of-type"
        table.querySelector(sel).click()
    })

}
