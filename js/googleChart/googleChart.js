google.load("visualization", "1.1", {packages:["bar"]});
google.setOnLoadCallback(drawStuff);

      function drawStuff() {
        
		var transferInData = new google.visualization.arrayToDataTable(transferIn);
        var options = {
          title: 'Transfer In',
          width: 900,
          legend: { position: 'none' },
          chart: { subtitle: 'Number Of Transfer In' },
          axes: {
            x: {
              0: { side: 'top', label: 'Transfer In'} // Top x-axis.
            }
          },
          bar: { groupWidth: "90%" }
        };

        var chart = new google.charts.Bar(document.getElementById('top_x_div'));
        // Convert the Classic options to Material options.
        chart.draw(transferInData, google.charts.Bar.convertOptions(options));
		
		//Transfer Out
		var transferOutData = new google.visualization.arrayToDataTable(transferOut);
		var optionsTransferOut = {
          title: 'Transfer Out',
          width: 900,
          legend: { position: 'none' },
          chart: { subtitle: 'Number Of Transfer Out' },
          axes: {
            x: {
              0: { side: 'top', label: 'Transfer Out'} // Top x-axis.
            }
          },
          bar: { groupWidth: "90%" }
        };

        var chartTransferOut = new google.charts.Bar(document.getElementById('top_transferOut_div'));
        // Convert the Classic options to Material options.
        chartTransferOut.draw(transferOutData, google.charts.Bar.convertOptions(optionsTransferOut));
      };