google.load("visualization", "1.1", {packages:["corechart"]});
google.setOnLoadCallback(drawStuff);

      function drawStuff() {
        
		var transferInData = new google.visualization.arrayToDataTable(transferIn);
        var transferInOptions = {
          title: 'Transfer In'
        };

        var transferInChart =  new google.visualization.PieChart(document.getElementById('top_transferIn_div'));
        // Convert the Classic options to Material options.
        transferInChart.draw(transferInData, transferInOptions);
		
		//Transfer Out
		var transferOutData = new google.visualization.arrayToDataTable(transferOut);
		var transferOutOptions = {
          title: 'Transfer Out'
        };

		var transferInChart =  new google.visualization.PieChart(document.getElementById('top_transferOut_div'));
        // Convert the Classic options to Material options.
        transferInChart.draw(transferOutData, transferOutOptions);
		
      };