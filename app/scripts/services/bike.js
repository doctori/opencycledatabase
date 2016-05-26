'use strict';

angular.module('bikeService', ['ngResource']).factory('Bike', function($resource){
	console.log($resource);
	return $resource('http://192.168.1.15:8080/bikes/:id',{id:'@ID'},{
		getRange: {
			url: 'http://192.168.1.15:8080/bikes?page=:page&per_page=:per_page&search_string=:search_string',
			method: 'GET',
			isArray:true,
			params: {
				page:'@page',
				per_page:'@per_page',
				search_string:'@search_string'
			}
		}
		//get:{isArray:true},
		// put:{method:'PUT', isArray:false},
		// post:{method:'POST', isArray:false},
		// update:{method: 'POST', isArray:false},

	});
}); 
