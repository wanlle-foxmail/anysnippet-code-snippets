def sort_tasks_by_dependency_order(tasks: dict[str, list[str]]) -> list[str]:
    """Return a stable task order that satisfies declared task dependencies."""
    # Flow:
    #   tasks -> validate input graph
    #            |
    #            +-> build in-degree and dependent edges
    #                 |
    #                 +-> queue tasks with no remaining prerequisites
    #                      |
    #                      +-> pop ready tasks and release dependents
    #                           |
    #                           +-> leftovers -> raise cycle error
    #                           +-> full order -> return ordered tasks
    if not isinstance(tasks, dict):
        raise TypeError("tasks must be a dictionary")

    in_degree: dict[str, int] = {}
    dependents: dict[str, list[str]] = {}

    for task_name, dependencies in tasks.items():
        if not isinstance(task_name, str):
            raise TypeError("task names must be strings")
        if not isinstance(dependencies, list):
            raise TypeError("dependencies must be lists of task names")

        seen_dependencies = set()
        for dependency_name in dependencies:
            if not isinstance(dependency_name, str):
                raise TypeError("dependency names must be strings")
            if dependency_name not in tasks:
                raise ValueError(f"unknown dependency: {dependency_name}")
            if dependency_name in seen_dependencies:
                raise ValueError(f"duplicate dependency: {dependency_name}")
            seen_dependencies.add(dependency_name)

        in_degree[task_name] = len(dependencies)
        dependents[task_name] = []

    for task_name, dependencies in tasks.items():
        for dependency_name in dependencies:
            dependents[dependency_name].append(task_name)

    ready_tasks = [task_name for task_name in tasks if in_degree[task_name] == 0]
    ordered_tasks = []

    while ready_tasks:
        task_name = ready_tasks.pop(0)
        ordered_tasks.append(task_name)

        for dependent_name in dependents[task_name]:
            in_degree[dependent_name] -= 1
            if in_degree[dependent_name] == 0:
                ready_tasks.append(dependent_name)

    if len(ordered_tasks) != len(tasks):
        raise ValueError("cyclic dependency detected")

    return ordered_tasks


if __name__ == "__main__":
    result = sort_tasks_by_dependency_order(
        {
            "package": ["test"],
            "test": ["build"],
            "build": [],
        }
    )
    print(result)