function barchartPlotter(e) {
    var ctx = e.drawingContext;
    var points = e.points;
    var y_bottom = e.dygraph.toDomYCoord(0);
    ctx.fillStyle = e.color;
    var min_sep = Infinity;

    for (var i = 1; i < points.length; i++) {
        var sep = points[i].canvasx - points[i - 1].canvasx;
        if (sep < min_sep) min_sep = sep;
    }

    var bar_width = Math.floor(2.0 / 3 * min_sep);
    for (var i = 0; i < points.length; i++) {
        var p = points[i];
        var center_x = p.canvasx;

        ctx.fillRect(center_x - bar_width / 2, p.canvasy,
            bar_width, y_bottom - p.canvasy);

        ctx.strokeRect(center_x - bar_width / 2, p.canvasy,
            bar_width, y_bottom - p.canvasy);
    }
}

(() => {
    
    var chartLayout = {
        showRangeSelector: true,
        legend: 'follow',
        fillGraph: true,
        colors: ['rgb(0,153,0)', 'orange', 'red'],
        stackedGraph: true,
        legendFormatter: agendasLegendFormatter,
        labelsSeparateLines: true,
        labelsKMB: true
    }

    function agendasLegendFormatter(data) {
        if (data.x == null) return '';
        var html = this.getLabels()[0] + ': ' + data.xHTML
        var total = data.series.reduce((total,n) => {
            return total + n.y
        }, 0)
        data.series.forEach((series) => {
            let percentage = ((series.y*100)/total).toFixed(2)
            html += `<br>${series.dashHTML}<span style="color: ${series.color};">${series.labelHTML}: ${series.yHTML} (${percentage}%)</span>`
        });
        return html
    }

    function cumulativeVoteChoicesData(d) {
        if (!d.yes instanceof Array) return [[0,0,0]]
        return d.yes.map((n,i) => {
            return [
                new Date(d.time[i]*1000),
                + d.yes[i],
                d.abstain[i], 
                d.no[i]
            ]
        })
    }

    function voteChoicesByBlockData(d) {
        if (!d.yes instanceof Array) return [[0,0,0,0]]
        return d.yes.map((n,i) => {
            return [
                d.height[i],
                + d.yes[i],
                d.abstain[i],
                d.no[i]
            ]
        });
    }

    function drawChart(el, data, options) {
        return new Dygraph(
            el,
            data,
            {
                ...chartLayout, 
                ...options
            }
        );
    }

    app.register("agenda", class extends Stimulus.Controller {
        static get targets() {
            return [
                "cumulativeVoteChoices",
                "voteChoicesByBlock"
            ]
        }

        connect() {
            var _this = this
            $.getScript("/js/dygraphs.min.js", function() {
                _this.drawCharts()
            })
        }

        disconnect(e) {
            this.cumulativeVoteChoicesChart.destroy()
            this.voteChoicesByBlockChart.destroy()   
        }

        drawCharts() {
            this.cumulativeVoteChoicesChart = drawChart(
                this.cumulativeVoteChoicesTarget,
                cumulativeVoteChoicesData(chartDataByTime),
                {
                    labels: ["Date", "Yes", "Abstain", "No"],
                    ylabel: 'Cumulative Vote Choices Cast',
                    title: 'Cumulative Vote Choices',
                    labelsKMB: true
                }
            );
            this.voteChoicesByBlockChart = drawChart(
                this.voteChoicesByBlockTarget,
                voteChoicesByBlockData(chartDataByBlock),
                {
                    labels: ["Block Height", "Yes", "Abstain", "No"],
                    ylabel: 'Vote Choices Cast',
                    title: 'Vote Choices By Block',
                    plotter: barchartPlotter
                }
            );
        }
    })
})()