#include <libdevmapper.h>
#include <lvm2app.h>
#include <stdio.h>

int wrapper_dm_list_iterate_items(struct dm_list *vgnames, char **r)
{
	struct lvm_str_list *strl;
	int num = 0;
	dm_list_iterate_items(strl, vgnames) {
		printf("--> %s\n", strl->str);
		r[num] = (char*)strl->str;
		num++;
	}
	return num;
}

// C helper functions:
char**makeCharArray(int size) {
        return calloc(sizeof(char*), size);
}

/*
static void setArrayString(char **a, char *s, int n) {
        a[n] = s;
}

static void freeCharArray(char **a, int size) {
        int i;
        for (i = 0; i < size; i++)
                free(a[i]);
        free(a);
}
*/

