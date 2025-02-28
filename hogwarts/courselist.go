//go:build !solution

package hogwarts

import "errors"

func GetCourseList(prereqs map[string][]string) []string {
	visited := make(map[string]int) // 0 - не visited; 1 - в процессе; 2 - visited
	courseList := []string{}
	for course := range prereqs {
		if visited[course] == 0 {
			courseList = DFS(course, prereqs, visited, courseList)
		}
	}
	return courseList
}

func DFS(course string, prereqs map[string][]string, visited map[string]int, courseList []string) []string {
	if visited[course] == 1 {
		panic(errors.New("Course " + course + " already visited"))
	}
	if visited[course] == 2 {
		return courseList
	}
	visited[course] = 1
	for _, prereq := range prereqs[course] {
		courseList = DFS(prereq, prereqs, visited, courseList)
	}
	visited[course] = 2
	courseList = append(courseList, course)
	return courseList
}
