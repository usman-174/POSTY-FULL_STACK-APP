import { HamburgerIcon, DeleteIcon,EditIcon } from '@chakra-ui/icons'
import { Menu, MenuButton, IconButton,Tooltip, MenuList, MenuItem,Portal } from '@chakra-ui/react'
import React from 'react'

const DropMenu : React.FC<{deletePost : Function}> = ({deletePost}): any=> {
return (
<Menu closeOnSelect={true} colorScheme="orange">
<Tooltip hasArrow label="Options."  bg="gray.400"aria-label="A tooltip">

  <MenuButton
    as={IconButton}
    aria-label="Options"
    icon={<HamburgerIcon />}
    variant="outline"
    />
    </Tooltip>
  <Portal>

  <MenuList bg="yellow.200" minWidth="11.5vw">
    <MenuItem icon={<DeleteIcon />} command="⌘D" onClick={()=>deletePost()}>
      Delete.
    </MenuItem>
    <MenuItem icon={<EditIcon />} command="⌘E">
      Edit.
    </MenuItem>
  </MenuList>
  </Portal>
</Menu>
)
}


export default DropMenu