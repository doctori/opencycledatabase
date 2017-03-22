'use strict';

angular.module('bikeService', ['ngResource']).factory('Bike', function($resource){
	console.log($resource);
	return $resource('http://127.0.0.1:8080/bikes/:id',{id:'@ID'},{
		// get:{isArray:true},
		// put:{method:'PUT', isArray:false},
		// post:{method:'POST', isArray:false},
		// update:{method: 'POST', isArray:false},

	});
}); 
