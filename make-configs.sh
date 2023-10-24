if [ ! -d "media" ]; then
    mkdir media
fi

if [ ! -d "env" ]; then
    mkdir env
    if [ ! -f "env/.env" ]; then
        cat << EOF > env/.env
PORT=8080

EMAIL=

FROM_EMAIL=

KEY=

PHONE=

EOF
    fi
fi
