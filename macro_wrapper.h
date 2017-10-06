#include <libdevmapper.h>
#include <lvm2app.h>
#include <stdio.h>

typedef char *names[];

struct result {
	char **namelist;
};

int wrapper_dm_list_iterate_items(struct dm_list *vgnames, char **r);
char**makeCharArray(int size);

