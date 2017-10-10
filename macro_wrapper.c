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

static int _set_pv_numeric_prop(pv_create_params_t pv_params, const char *name,
                                                                unsigned long long value)
{
        struct lvm_property_value prop_value = {
                .is_integer = 1,
                .value.integer = value,
        };

        return lvm_pv_params_set_property(pv_params, name, &prop_value);
}

#define SET_PV_PROP(params, name, value) \
        do { \
                if (_set_pv_numeric_prop(params, name, value) == -1) \
                        goto error; \
        } while(0)\


void wrapper_set_pv_prop(pv_create_params_t pv_params, char *name, long long value) {
	SET_PV_PROP(pv_params, name, value);
error:
	// TODO
	return;
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


int getN(void* p) {
	lvm_property_value_t t = *(lvm_property_value_t *)p;
	return t.is_valid;
}
