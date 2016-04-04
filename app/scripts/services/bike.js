angular.module('bikeService', ['ngResource']).factory('Bike', function($resource){
	"use strict";
	console.log($resource);
	return $resource('http://192.168.198.130:8080/bikes/:id',{id:'@ID'},{
		// get:{isArray:true},
		// put:{method:'PUT', isArray:false},
		// post:{method:'POST', isArray:false},
		// update:{method: 'POST', isArray:false},

	});
}); 