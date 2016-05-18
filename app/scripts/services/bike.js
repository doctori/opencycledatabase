'use strict';

angular.module('bikeService', ['ngResource']).factory('Bike', function($resource){
	console.log($resource);
	return $resource('http://192.168.1.15:8080/bikes/:id?per_page=50',{id:'@ID'},{
		// get:{isArray:true},
		// put:{method:'PUT', isArray:false},
		// post:{method:'POST', isArray:false},
		// update:{method: 'POST', isArray:false},

	});
}); 
