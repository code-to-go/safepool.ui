<script>
  import { BlockTitle, f7, Icon, Link, List, ListItem, Navbar, NavLeft, NavRight, NavTitle, Page, Searchbar, Subnavbar, theme } from 'framework7-svelte';
  import { onMount } from 'svelte'
  import { GetPoolList, ListLibrary } from '../../wailsjs/go/main/App';
  
  export let poolName;
  let vlData = {
    items: [],
  }
  onMount( async () => {
    const documents = await ListLibrary(poolName) || []
    console.log(documents)
    vlData.items = documents.map(d => ({
      title: `${d.name}`,
      subtitle: `${d.localPath}`,
    }))
  })

  function searchAll(query, items) {
    const found = [];
    for (let i = 0; i < items.length; i += 1) {
      if (items[i].title.toLowerCase().indexOf(query.toLowerCase()) >= 0 || query.trim() === '') found.push(i);
    }
    return found; // return array with mathced indexes
  }

  function renderExternal(virtualList, virtualListData) {
    vlData = virtualListData;
  }



</script>

<Page>
  <Navbar title="Virtual List">
    <Subnavbar inner={false}>
      <Searchbar
        searchContainer=".virtual-list"
        searchItem="li"
        searchIn=".item-title"
        disableButton={!theme.aurora}
      ></Searchbar>
    </Subnavbar>
  </Navbar>
  <List class="searchbar-not-found">
    <ListItem title="Nothing found"></ListItem>
  </List>
  <List
    class="searchbar-found"
    ul={false}
    medialList
    virtualList
    virtualListParams={{
      items: vlData.items,
      searchAll: searchAll,
      renderExternal: renderExternal,
      height: theme.ios ? 63 : (theme.md ? 73 : 77)
    }}
  >
    <ul>
      {#each vlData.items as item, index (index)}
        <ListItem
          mediaItem
          link="#"
          title={item.title}
          subtitle={item.subtitle}
          style={`top: ${vlData.topPosition}px`}
          virtualListIndex={vlData.items.indexOf(item)}
        ></ListItem>
      {/each}
    </ul>
  </List>
</Page>
  


<style>
</style>
